// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.messages;

import java.io.IOException;
import java.nio.ByteBuffer;
import java.util.Arrays;

public class ReadOnlyConsoleMessage extends ConsoleMessage {

    public final String source;
    public final String[] text;

    public ReadOnlyConsoleMessage(
        final String source,
        final String[] text
    ) {
        this.source = source;
        this.text = text;
    }

    public static ReadOnlyConsoleMessage deserialize(ByteBuffer buffer) throws IOException {
        var source = deserializeString(buffer);
        var text = deserializeStringArray(buffer);
        return new ReadOnlyConsoleMessage(source, text);
    }

    public byte[] serialize() throws IOException {
        var buffer = ByteBuffer.allocate(1024);
        super.serializePartial(buffer);
        buffer.putInt(Type.CONSOLE_READ_ONLY.value);
        serializeString(buffer, this.source);
        serializeStringArray(buffer, this.text);
        serializeFixLength(buffer);
        return Arrays.copyOf(buffer.array(), buffer.position());
    }

    @Override
    public String toString() {
        var sb = new StringBuilder();
        sb.append(">>");
        if (this.source.length() > 0) {
            sb.append(this.source).append("*");
        }

        sb.append(this.text[0]);
        return sb.toString();
    }
}
