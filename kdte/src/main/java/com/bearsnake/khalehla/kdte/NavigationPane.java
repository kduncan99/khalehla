package com.bearsnake.khalehla.kdte;

import javafx.scene.control.TreeItem;
import javafx.scene.control.TreeView;
import javafx.scene.layout.VBox;

public class NavigationPane extends VBox /*TreeView<String>*/ {

    private TreeView<String> createDirectoriesView() {
        var root = new TreeItem<>("Directories");
        var item = new TreeItem<>("local");
        root.getChildren().add(item);
        root.setExpanded(true);
        // TODO create other configured directory entries
        return new TreeView<>(root);
    }

    private TreeView<String> createMediaPoolsView() {
        var root = new TreeItem<>("Media Pools");
        // TODO create other configured media pool entries
        root.setExpanded(true);
        return new TreeView<>(root);
    }

    private TreeView<String> createSystemsView() {
        var root = new TreeItem<>("Systems");
        var item = new TreeItem<>("local");
        root.getChildren().add(item);
        // TODO create other configured system entries
        root.setExpanded(true);
        return new TreeView<>(root);
    }

    public NavigationPane() {
        this.getChildren().add(createSystemsView());
        this.getChildren().add(createDirectoriesView());
        this.getChildren().add(createMediaPoolsView());
    }
}
