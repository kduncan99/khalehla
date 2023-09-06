// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.console;

import javafx.scene.layout.VBox;

import java.io.IOException;

public class Console extends VBox {

    private final String title;
    private final Transport transport;

    public Console(String title, String endPoint, int port) {
        this.title = title;
        this.transport = new Transport(endPoint, port);
    }

    public void connect() throws IOException {
        this.transport.connect();
    }

    public byte[] readInput() throws IOException {
        return this.transport.readInput();
    }
}
