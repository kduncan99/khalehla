// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved
//
// Khalehla DeskTop Environment

package com.bearsnake.khalehla.kdte;

import com.bearsnake.khalehla.kdte.console.Console;
import com.bearsnake.khalehla.kdte.console.ConsolePane;
import javafx.application.Application;
import javafx.scene.Scene;
import javafx.scene.control.Menu;
import javafx.scene.control.MenuBar;
import javafx.scene.control.MenuItem;
import javafx.scene.layout.BorderPane;
import javafx.stage.Stage;

import java.io.IOException;

/*
  --------------------------------------------------------------------------------------
  | menu                                                                               |
  |------------------------------------------------------------------------------------|
  | Systems          | local-MFD  local-DEMAND  local-TIP  local-CONSOLE  tip1-CONSOLE |
  |   local          | production-DEMAND  local-FILES  local#SYS$*INFO$(32)            |
  |   production     | local#BOB*BOB.ACCT-REPORT/TEST(COB) local-MEDIA.SYS001          |
  |   prodbackup     | production-QUEUES:PR1                                           |
  |   tip1           |-----------------------------------------------------------------|
  |   tip2           |                                                                 |
  |   production     |                                                                 |
  |                  |                                                                 |
  | Media Pools      |                                                                 |
  |   local          |                                                                 |
  |   software       |                                                                 |
  |   backups        |                                                                 |
  |                  |                                                                 |
  --------------------------------------------------------------------------------------
 */
public class Main extends Application {

    public static final String TITLE = "Khalehla DeskTop Environment";
    public static final String VERSION = "1.0";

    public static MenuBar createMenu() {
        var menuApplicationAbout = new MenuItem("About");
        var menuApplicationQuit = new MenuItem("Quit");
        var menuApplication = new Menu("Application");
        menuApplication.getItems().addAll(menuApplicationAbout, menuApplicationQuit);
        var menuFile = new Menu("File");
        var menuEdit = new Menu("Edit");
        var menuSystem = new Menu("System");
        var menuHelp = new Menu("Help");

        var menuBar = new MenuBar();
        menuBar.getMenus().clear();
        menuBar.getMenus().addAll(menuApplication, menuFile, menuEdit, menuSystem, menuHelp);
        return menuBar;
    }

    @Override
    public void start(Stage stage) throws IOException {
        // TODO probably should move this root crap into a separate class
        BorderPane root = new BorderPane();

        root.setTop(createMenu());
        root.setLeft(new NavigationPane());
        root.setCenter(new ContentPane());

        Scene scene = new Scene(root);
        stage.setTitle(TITLE + " - " + VERSION);
        stage.setScene(scene);
        stage.show();
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        launch();
//        var c = new Console("Dork", "127.0.0.1", 2200);
//        c.connect();
//
//        var pendingMessage = ByteBuffer.allocate(1024);
//        while (true) {
//            var in = c.readInput();
//            if (in.length > 0) {
//                pendingMessage.put(in);
//                if (pendingMessage.position() > 8) {
//                    var slice = pendingMessage.slice(4, 4);
//                    var msgLen = slice.getInt();
//                    if (msgLen <= pendingMessage.position()) {
//                        slice = pendingMessage.slice(0, msgLen);
//                        var msg = Message.deserialize(slice);
//                        System.out.println(msg.toString());
//
//                        var remainingArray = Arrays.copyOfRange(pendingMessage.array(), msgLen, pendingMessage.position());
//                        pendingMessage.clear();
//                        pendingMessage.put(remainingArray);
//                    }
//                }
//            } else {
//                try {
//                    Thread.sleep(100);
//                } catch (InterruptedException ex) {
//                    // nothing to be done
//                }
//            }
//        }
    }
}
