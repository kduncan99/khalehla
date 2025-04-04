// Khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.console;

import javafx.scene.layout.Pane;

import java.io.IOException;

public class Console {

    private final Transport transport;
    private final ConsolePane pane;

    public Console(String title, String endPoint, int port) {
        this.transport = new Transport(endPoint, port);
        this.pane = new ConsolePane();
    }

    public void connect() throws IOException {
        this.transport.connect();
    }

    public Pane getPane() { return this.pane; }

    public byte[] readInput() throws IOException {
        return this.transport.readInput();
    }
}
