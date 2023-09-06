// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.messages;

import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.HashMap;

// Serialized form of message:
//  uint32 identifier:
//      universal value for all messages: 2200
//  uint32 length:
//      length of the entire message - when invoked, the given buffer must be this length
//  uint32 category:
//      indicates the message category
//      categories are defined by the various packages which utilize this functionality.
//      each package corresponds to a category, and all category values are necessarily unique.
//  uint32 type:
//      indicates the message type
//      types are defined by the package to which the message belongs, and must be unique within the package category.
//  byte[] payload:
//      defined by the category and the particular message type within the category.
public abstract class Message {

    public static final int IDENTIFIER = 2200;

    public abstract static class Deserializer {
        public abstract Message deserialize(ByteBuffer buffer, Category category, Type type) throws IOException;
    }

    private static final HashMap<Category, Deserializer> deserializers = new HashMap<>();
    static {
        deserializers.put(Category.SYSTEM, new SystemMessage.SystemMessageDeserializer());
        deserializers.put(Category.CONSOLE, new ConsoleMessage.ConsoleMessageDeserializer());
    }

    public static void registerDeserializer(Category category, Deserializer d) {
        synchronized (deserializers) {
            deserializers.put(category, d);
        }
    }

    public void serializePartial(ByteBuffer buffer) throws IOException {
        buffer.putInt(IDENTIFIER);
        buffer.putInt(0); // length - set to zero for now, and we'll fix it later.
    }

    public final void serializeFixLength(ByteBuffer buffer) throws IOException {
        buffer.putInt(1, buffer.position());
    }

    public static Message deserialize(ByteBuffer buffer) throws IOException {
        buffer.position(0);

        var ident = buffer.getInt();
        if (ident != IDENTIFIER) {
            throw new IOException("Buffer does not contain a Khalehla network message");
        }

        buffer.getInt(); // skip length - we don't need it anymore

        Category category = Category.getCategory(buffer.getInt());
        var d = deserializers.get(category);
        if (d == null) {
            throw new IOException(String.format("No deserializer for message category %d", category.value));
        }

        Type type = Type.getType(buffer.getInt());
        return d.deserialize(buffer, category, type);
    }

    protected static String deserializeString(ByteBuffer buffer) throws IOException {
        var length = buffer.getInt();
        if (length > buffer.remaining()) {
            throw new IOException("Cannot deserialize string - out of data");
        }

        var limit = buffer.position() + length;
        var str = new String(Arrays.copyOfRange(buffer.array(), buffer.position(), limit), StandardCharsets.UTF_8);
        buffer.position(limit);
        return str;
    }

    protected static String[] deserializeStringArray(ByteBuffer buffer) throws IOException {
        var count = buffer.getInt();
        var result = new String[count];
        for (var x = 0; x < count; x++) {
            result[x] = deserializeString(buffer);
        }
        return result;
    }

    protected static void serializeString(ByteBuffer buffer, String str) throws IOException {
        buffer.putInt(str.length());
        buffer.put(str.getBytes(StandardCharsets.UTF_8));
    }

    protected static void serializeStringArray(ByteBuffer buffer, String[] array) throws IOException {
        buffer.putInt(array.length);
        for (var str : array) {
            serializeString(buffer, str);
        }
    }
}
