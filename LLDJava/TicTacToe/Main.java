package LLDJava.TicTacToe;

public class Main {
    public static void main(String[] args) {
        Game game = new Game();
        String result = game.startGame();
        if(result == "Tie"){
            System.out.println("Game is " +  result);
        }else{
            System.out.println("Winner is " + result);
        }
    }
}
