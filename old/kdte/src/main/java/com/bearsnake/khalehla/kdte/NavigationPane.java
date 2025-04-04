package com.bearsnake.khalehla.kdte;

import javafx.scene.control.TreeItem;
import javafx.scene.control.TreeView;

public class NavigationPane extends TreeView<String> {

    private final TreeItem<String> systemsMainItem = new TreeItem<>("Systems");
    private final TreeItem<String> directoriesMainItem = new TreeItem<>("Directories");
    private final TreeItem<String> mediaPoolsMainItem = new TreeItem<>("Media Pools");
    private final TreeItem<String> rootItem = new TreeItem<>("Resources");

    //  must be called on the graphics thread
    public void populateFromConfig(/* TODO resource list */) {
        this.systemsMainItem.getChildren().clear();
        this.directoriesMainItem.getChildren().clear();
        this.mediaPoolsMainItem.getChildren().clear();

        // TODO temporary code
        this.systemsMainItem.getChildren().add(new TreeItem<>("local"));
        this.systemsMainItem.getChildren().add(new TreeItem<>("Dev"));
        this.systemsMainItem.getChildren().add(new TreeItem<>("Production"));
        this.systemsMainItem.getChildren().add(new TreeItem<>("TIP Cloud"));
        this.systemsMainItem.setExpanded(true);

        this.directoriesMainItem.getChildren().add(new TreeItem<>("local"));
        this.directoriesMainItem.getChildren().add(new TreeItem<>("Dev"));
        this.directoriesMainItem.setExpanded(true);

        this.mediaPoolsMainItem.getChildren().add(new TreeItem<>("Dev"));
        this.mediaPoolsMainItem.getChildren().add(new TreeItem<>("Production"));
        this.mediaPoolsMainItem.getChildren().add(new TreeItem<>("Software"));
        this.mediaPoolsMainItem.setExpanded(true);
        // TODO end temporary code
    }

    public NavigationPane() {
        populateFromConfig();

        this.rootItem.getChildren().add(this.systemsMainItem);
        this.rootItem.getChildren().add(this.directoriesMainItem);
        this.rootItem.getChildren().add(this.mediaPoolsMainItem);
        this.rootItem.setExpanded(true);

        this.setRoot(this.rootItem);
    }
}
