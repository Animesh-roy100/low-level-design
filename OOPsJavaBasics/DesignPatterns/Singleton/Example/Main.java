package OOPsJavaBasics.DesignPatterns.Singleton.Example;

public class Main {
    public static void main(String[] args) {
        Settings setting = Settings.getObject();
		System.out.println("Setting Object Color: " + setting.getColor());
		System.out.println("Font: " + setting.getFont());
		Settings set1=Settings.getObject();
		set1.setColor("Red");
		System.out.println("Set Object Color: " + set1.getColor());
		System.out.println("Font: " + set1.getFont());
		System.out.println("Setting Object Color: " + setting.getColor());
    }
}
