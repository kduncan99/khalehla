// Khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte;

import com.bearsnake.khalehla.kdte.console.Console;
import javafx.geometry.Insets;
import javafx.geometry.Pos;
import javafx.scene.control.Tab;
import javafx.scene.control.TabPane;
import javafx.scene.layout.*;
import javafx.scene.paint.Color;

// Container class for any specific content, especially fixed-size content which wants to be
// vertically and horizontally centered when the overall application window expands.
public class ContentPane extends TabPane {

    public ContentPane() {
        //  TODO temporary code
        var console = new Console("CONSL0", "127.0.0.1", 2200);

        var container = new GridPane();
        container.setAlignment(Pos.CENTER);
        container.getChildren().add(console.getPane());

        var tab = new Tab();
        tab.setText("local:console");
        tab.setContent(container);

        getTabs().add(tab);
        setBackground(new Background(new BackgroundFill(Color.DARKGRAY, CornerRadii.EMPTY, Insets.EMPTY)));
    }
}
