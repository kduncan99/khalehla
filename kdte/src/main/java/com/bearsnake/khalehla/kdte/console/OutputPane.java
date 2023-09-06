// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.console;

import javafx.scene.layout.Pane;
import javafx.scene.paint.Color;
import javafx.scene.text.Font;

public class OutputPane extends Pane {

    private final static int COLUMN = 80;
    private final static int ROWS = 24;

    private final static Color DEFAULT_INPUT_COLOR = Color.WHITE;
    private final static Color DEFAULT_READ_ONLY_COLOR = Color.GREEN;
    private final static Color DEFAULT_READ_REPLY_COLOR = Color.RED;
    private final static Color DEFAULT_STATUS_COLOR = Color.CYAN;

    private Font font;
    private Color[] inputColor = {DEFAULT_INPUT_COLOR, Color.BLACK};
    private Color[] readOnlyColor = {DEFAULT_READ_ONLY_COLOR, Color.BLACK};
    private Color[] readReplyColor = {DEFAULT_READ_REPLY_COLOR, Color.BLACK};
    private Color[] statusColor = {DEFAULT_STATUS_COLOR, Color.BLACK};

    public OutputPane() {
        this.font = Font.font("Courier New", 20);
    }
}
