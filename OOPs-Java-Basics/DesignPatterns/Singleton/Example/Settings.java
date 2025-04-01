package oopsJava.DesignPatterns.Singleton.Example;

public class Settings {
    private static Settings firstObject = null;
    private String font = "calibri";
    private String color = "black";

    private Settings(){}

    public static Settings getObject() {
        if(firstObject == null) {
            firstObject = new Settings();
        }

        return firstObject;
    }

    public void setFont(String font) {
        this.font = font;
    }

    public void setColor(String color) {
        this.color = color;
    }

    public String getFont() {
        return this.font;
    }

    public String getColor() {
        return this.color;
    }
}
