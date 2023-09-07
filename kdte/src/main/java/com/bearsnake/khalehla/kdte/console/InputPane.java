// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.console;

import javafx.scene.canvas.Canvas;
import javafx.scene.paint.Color;
import javafx.scene.text.Font;
import javafx.scene.text.TextAlignment;

public class InputPane extends InputOutputPane {

    private Canvas canvas;

    public InputPane(
        Font font,
        double width,
        double height
    ) {
        super(font, width, height);
    }

    void putText(
        final Color fgColor,
        final Color bgColor,
        final String text
    ) {
        this.canvas = new Canvas();
        this.canvas.setHeight(this.height);
        this.canvas.setWidth(this.width);

        var gc = this.canvas.getGraphicsContext2D();
//        gc.setStroke(bgColor);
        gc.setFill(bgColor);
        gc.setFont(this.font);
        gc.fillRect(0, 0, width, height);
//        gc.strokeRect(0, 0, width, height);

        gc.setTextAlign(TextAlignment.LEFT);
        gc.setStroke(fgColor);
        gc.strokeText(text, 0, this.height - 3);
        gc.closePath();
        getChildren().clear();
        getChildren().add(this.canvas);
        this.updateBounds();//TODO what?
    }
}
