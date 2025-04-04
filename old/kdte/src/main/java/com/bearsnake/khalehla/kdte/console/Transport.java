// Khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte.console;

import javax.net.SocketFactory;
import javax.net.ssl.SSLSocket;
import javax.net.ssl.SSLSocketFactory;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;

public class Transport {

    public final String endPoint;
    public final int port;

    private Socket socket;
    private InputStream inputStream;
    private OutputStream outputStream;

    public Transport(String endPoint, int port) {
        this.endPoint = endPoint;
        this.port = port;

        this.socket = null;
        this.inputStream = null;
        this.outputStream = null;
    }

    public void connect() throws IOException {
        var sf = SocketFactory.getDefault();
        this.socket = sf.createSocket(this.endPoint, this.port);
        this.inputStream = this.socket.getInputStream();
        this.outputStream = this.socket.getOutputStream();
    }

    public void secureConnect() throws IOException {
        var ssf = SSLSocketFactory.getDefault();
        var ss = (SSLSocket) ssf.createSocket(this.endPoint, this.port);
        ss.setUseClientMode(true);
        ss.startHandshake();
        this.socket = ss;
        this.inputStream = this.socket.getInputStream();
        this.outputStream = this.socket.getOutputStream();
    }

    public void disconnect() throws IOException {
        if (this.socket != null) {
            this.inputStream.close();
            this.outputStream.close();
            this.socket.close();
            this.socket = null;
        }
    }

    public byte[] readInput() throws IOException {
        var av = this.inputStream.available();
        if (av > 0) {
            return this.inputStream.readNBytes(av);
        } else {
            return new byte[0];
        }
    }
}
