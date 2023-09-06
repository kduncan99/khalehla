module com.bearsnake.khalehla.kdte {
    requires javafx.controls;
    requires javafx.fxml;

    requires org.controlsfx.controls;

    opens com.bearsnake.khalehla.kdte to javafx.fxml;
    exports com.bearsnake.khalehla.kdte;
    exports com.bearsnake.khalehla.kdte.messages;
    opens com.bearsnake.khalehla.kdte.messages to javafx.fxml;
}