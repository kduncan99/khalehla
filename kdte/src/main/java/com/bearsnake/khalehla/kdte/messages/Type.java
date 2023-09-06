package com.bearsnake.khalehla.kdte.messages;

import java.util.Arrays;

public enum Type {
    SYSTEM_HEARTBEAT(1),

    CONSOLE_READ_ONLY(101),
    CONSOLE_READ_REPLY(102),
    CONSOLE_RESET(103),
    CONSOLE_STATUS(104),
    CONSOLE_SOLICITED(105),
    CONSOLE_UNSOLICITED(106);

    public final int value;

    Type(int value) {
        this.value = value;
    }

    public static Type getType(int value) {
        return Arrays.stream(Type.values()).filter(t -> value == t.value).findFirst().orElse(null);
    }
}
