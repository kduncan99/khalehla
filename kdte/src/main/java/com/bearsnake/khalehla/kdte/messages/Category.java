package com.bearsnake.khalehla.kdte.messages;

import java.util.Arrays;

public enum Category {
    SYSTEM(0),
    CONSOLE(1);

    public final int value;

    Category(int value) {
        this.value = value;
    }

    public static Category getCategory(int value) {
        return Arrays.stream(Category.values()).filter(c -> value == c.value).findFirst().orElse(null);
    }
}
