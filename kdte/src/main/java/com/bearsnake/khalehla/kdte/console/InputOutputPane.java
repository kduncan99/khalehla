// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.console;

import javafx.scene.layout.Pane;
import javafx.scene.text.Font;

public abstract class InputOutputPane extends Pane {

    protected Font font;
    protected double width;
    protected double height;

    public InputOutputPane(
        Font font,
        double width,
        double height
    ) {
        this.font = font;
        setSize(width, height);
    }

    void setFont(
        Font font
    ) {
        this.font = font;
    }

    void setSize(
        double width,
        double height
    ) {
        this.width = width;
        this.height = height;

        this.setMaxSize(this.width, this.height);
        this.setMinSize(this.width, this.height);
        this.setPrefSize(this.width, this.height);
    }
}
