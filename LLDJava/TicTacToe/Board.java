package LLDJava.TicTacToe;

public class Board {
    public int size;
    public Cell[][] board;

    Board(int size) {
        this.size = size;
        board = new Cell[size][size];
    }

    boolean getFreeCells(){
        for(int i = 0;i<size;i++){
            for(int j = 0;j<size;j++){
                if(board[i][j] == null){
                    return true;
                }
            }
        }
        return false;
    }

    boolean addAtPosition(int x,int y, Cell cell){
        if(board[x][y] != null){
            return false;
        }
        board[x][y] = cell;
        return true;
    }
    
    public void printBoard() {
        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                if (board[i][j] != null) {
                    System.out.print(board[i][j].cellType.name() + "   ");
                } else {
                    System.out.print("    ");
                }
                System.out.print(" | ");
            }
            System.out.println();
        }
    }
}
