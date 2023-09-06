// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.messages;

import java.io.IOException;
import java.nio.ByteBuffer;

public abstract class SystemMessage extends Message {

    protected static class SystemMessageDeserializer extends Message.Deserializer {
        @Override
        public Message deserialize(ByteBuffer buffer, Category category, Type type) throws IOException {
            buffer.putInt(Category.SYSTEM.value);
            return switch (type) {
                case SYSTEM_HEARTBEAT -> HeartbeatSystemMessage.deserialize(buffer);
                default -> null;
            };
        }
    }

    @Override
    public void serializePartial(ByteBuffer buffer) throws IOException {
        super.serializePartial(buffer);
        buffer.putInt(Category.SYSTEM.value);
    }
}
