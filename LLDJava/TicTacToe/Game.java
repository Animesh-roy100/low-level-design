package LLDJava.TicTacToe;

import java.util.Deque;
import java.util.LinkedList;
import java.util.Scanner;

public class Game {
    Deque<Player> players;
    Board gameBoard;

    Game(){
        players = new LinkedList<Player>();
        CellX cellX = new CellX();
        players.add(new Player("Player1",cellX));
        CellO cellY = new CellO();
        players.add((new Player("Player2",cellY)));
        gameBoard = new Board(3);
    }

    public String startGame(){
        boolean isWinner = true;
        while(isWinner){
            Player playerturn = players.removeFirst();
            gameBoard.printBoard();
            boolean freeCells = gameBoard.getFreeCells();
            if(!freeCells){
                isWinner = false;
                continue;
            }

            System.out.print("Player " + playerturn.getName() + "Enter Position");
            try (Scanner inputScanner = new Scanner(System.in)) {
                String s = inputScanner.nextLine();
                String[] values = s.split(",");
                int inputRow = Integer.valueOf(values[0]);
                int inputColumn = Integer.valueOf(values[1]);

                boolean isValidPosition = gameBoard.addAtPosition(inputRow,inputColumn,playerturn.cellType);

                if(!isValidPosition){
                    System.out.print("Please Enter valid Position");
                    players.addFirst(playerturn);
                    continue;
                }
                players.addLast(playerturn);

                boolean iswinner = isWinner(inputRow,inputColumn,playerturn.cellType.cellType);
                if(iswinner){
                     return playerturn.name;
                }
            } catch (NumberFormatException e) {
                // TODO Auto-generated catch block
                e.printStackTrace();
            }
        }
        return "Match Tie";
    }

    public boolean isWinner(int row,int column,CellType cellType){
        boolean isrow = true;
        boolean isColumn = true;
        boolean isDiagonal = true;
        boolean isAntiDiagonal = true;
        for(int i = 0;i<gameBoard.size;i++){
            if(gameBoard.board[row][i] == null || gameBoard.board[row][i].cellType != cellType){
                isrow = false;
            }
        }
        for(int i = 0;i<gameBoard.size;i++){
            if(gameBoard.board[i][column] == null || gameBoard.board[i][column].cellType != cellType){
                isColumn = false;
            }
        }
        for(int i = 0,j=0;i<gameBoard.size;i++,j++){
            if(gameBoard.board[i][j] == null || gameBoard.board[i][j].cellType != cellType){
                isDiagonal = false;
            }
        }
        for(int i = 0,j= gameBoard.size-1;i< gameBoard.size;i++,j--){
            if(gameBoard.board[i][j] == null || gameBoard.board[i][j].cellType != cellType){
                isAntiDiagonal = false;
            }
        }
        return isrow || isColumn || isDiagonal || isAntiDiagonal;
    }
}
