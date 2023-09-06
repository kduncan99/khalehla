// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.messages;

import java.io.IOException;
import java.nio.ByteBuffer;

public class ConsoleMessage extends Message {

    protected static class ConsoleMessageDeserializer extends Message.Deserializer {
        @Override
        public Message deserialize(ByteBuffer buffer, Category category, Type type) throws IOException {
            return switch (type) {
                case CONSOLE_READ_ONLY -> ReadOnlyConsoleMessage.deserialize(buffer);
                case CONSOLE_READ_REPLY -> null; // TODO
                case CONSOLE_RESET -> null; // TODO
                case CONSOLE_STATUS -> null; // TODO
                case CONSOLE_SOLICITED -> null; // TODO
                case CONSOLE_UNSOLICITED -> null; // TODO
                default -> null;
            };
        }
    }

    @Override
    public void serializePartial(ByteBuffer buffer) throws IOException {
        super.serializePartial(buffer);
        buffer.putInt(Category.CONSOLE.value);
    }
}
