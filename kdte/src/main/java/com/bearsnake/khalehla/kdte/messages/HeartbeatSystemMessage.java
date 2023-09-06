// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.messages;

import java.io.IOException;
import java.nio.ByteBuffer;
import java.util.Arrays;

public class HeartbeatSystemMessage extends SystemMessage {

    protected HeartbeatSystemMessage() {}

    public static Message deserialize(ByteBuffer buffer) throws IOException {
        return new HeartbeatSystemMessage();
    }

    public byte[] serialize() throws IOException {
        var buffer = ByteBuffer.allocate(1024);
        super.serializePartial(buffer);
        buffer.putInt(Type.SYSTEM_HEARTBEAT.value);
        serializeFixLength(buffer);
        return Arrays.copyOf(buffer.array(), buffer.position());
    }
}
