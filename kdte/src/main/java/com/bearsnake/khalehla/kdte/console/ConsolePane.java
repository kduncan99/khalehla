// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.console;

import javafx.geometry.Bounds;
import javafx.scene.layout.*;
import javafx.scene.paint.Color;
import javafx.scene.text.Font;
import javafx.scene.text.Text;

public class ConsolePane extends VBox {

    protected final static int COLUMNS = 80;
    protected final static int ROWS = 24;
    private final static double BORDER_WIDTH = 5.0;
    private final static double VERTICAL_SPACING = 5.0;
    private static final int DEFAULT_FONT_SIZE = 12;
    private final static Color DEFAULT_INPUT_COLOR = Color.WHITE;
    private final static Color DEFAULT_READ_ONLY_COLOR = Color.GREEN;
    private final static Color DEFAULT_READ_REPLY_COLOR = Color.RED;
    private final static Color DEFAULT_STATUS_COLOR = Color.CYAN;

    private final Font displayFont;
    private final OutputPane outputPane;
    private final InputPane inputPane;

    private static final String BLANK_TEXT_FORMATTER = String.format("%%%ds", COLUMNS);
    protected static final String BLANK_TEXT = String.format(BLANK_TEXT_FORMATTER, "X");

    private Color[] inputColors = {DEFAULT_INPUT_COLOR, Color.BLACK};
    private Color[] readOnlyColors = {DEFAULT_READ_ONLY_COLOR, Color.BLACK};
    private Color[] readReplyColors = {DEFAULT_READ_REPLY_COLOR, Color.BLACK};
    private Color[] statusColors = {DEFAULT_STATUS_COLOR, Color.BLACK};

    public ConsolePane() {
        this.displayFont = new Font("Courier New", DEFAULT_FONT_SIZE);

        var text = new Text(BLANK_TEXT);
        text.setFont(this.displayFont);
        var bounds = text.getLayoutBounds();

        this.outputPane = new OutputPane(this.displayFont, bounds.getWidth(), ROWS * bounds.getHeight());
        this.outputPane.setBackground(new Background(new BackgroundFill(Color.BLUE, null, null)));
        this.inputPane = new InputPane(this.displayFont, bounds.getWidth(), bounds.getHeight());
        this.inputPane.setBackground(new Background(new BackgroundFill(Color.GREEN, null, null)));

        getChildren().add(this.outputPane);
        getChildren().add(this.inputPane);
        setSpacing(5.0);
        setBorder(new Border(new BorderStroke(Color.LIGHTGRAY,
                                              BorderStrokeStyle.SOLID,
                                              CornerRadii.EMPTY,
                                              new BorderWidths(BORDER_WIDTH))));
        setBackground(new Background(new BackgroundFill(Color.LIGHTGRAY, null, null)));

        var totalWidth = bounds.getWidth() + 2 * BORDER_WIDTH;
        var totalHeight = BORDER_WIDTH
                          + ROWS * bounds.getHeight()
                          + VERTICAL_SPACING
                          + bounds.getHeight()
                          + BORDER_WIDTH;
        System.out.printf("%f x %f\n", totalWidth, totalHeight);//TODO remove
        this.setPrefWidth(totalWidth);
        this.setMaxWidth(totalWidth);
        this.setMinWidth(totalWidth);
        this.setPrefHeight(totalHeight);
        this.setMaxHeight(totalHeight);
        this.setMinHeight(totalHeight);

        this.inputPane.putText(Color.WHITE, Color.BLACK, "INPUT AREA HERE");
    }

    protected static Bounds getInputOutputPaneBounds(Font font) {
        var text = new Text(BLANK_TEXT);
        text.setFont(font);
        return text.getLayoutBounds();
    }
}
