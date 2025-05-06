package LLDJava.SnakeLadder;

import java.util.concurrent.ThreadLocalRandom;

public class Board {
    Cell[][] cells;

    public Board(int boardSize, int numberOfSnakes, int numberOfLadders) {
        initializeCells(boardSize);
        addSnakesAndLadders(numberOfSnakes, numberOfLadders);
    }

    private void initializeCells(int boardSize) {
        cells = new Cell[boardSize][boardSize];
        for (int i = 0; i < boardSize; i++) {
            for (int j = 0; j < boardSize; j++) {
                cells[i][j] = new Cell();
            }
        }
    }

    private void addSnakesAndLadders(int numberOfSnakes, int numberOfLadders) {
        while(numberOfSnakes > 0) {
           int snakeHead = ThreadLocalRandom.current().nextInt(1,cells.length*cells.length-1);
           int snakeTail = ThreadLocalRandom.current().nextInt(1,cells.length*cells.length-1);
           if(snakeTail >= snakeHead) {
               continue;
           }

           Jump snakeObj = new Jump();
           snakeObj.startPosition = snakeHead;
           snakeObj.endPosition = snakeTail;

           Cell cell = getCell(snakeHead);
           cell.jump = snakeObj;

            numberOfSnakes--;
        }

        while(numberOfLadders > 0) {
            int ladderStart = ThreadLocalRandom.current().nextInt(1,cells.length*cells.length-1);
            int ladderEnd = ThreadLocalRandom.current().nextInt(1,cells.length*cells.length-1);
            if(ladderStart >= ladderEnd) {
                continue;
            }

            Jump ladderObj = new Jump();
            ladderObj.startPosition = ladderStart;
            ladderObj.endPosition = ladderEnd;

            Cell cell = getCell(ladderStart);
            cell.jump = ladderObj;

            numberOfLadders--;
        }
    }

    Cell getCell(int playerPosition) {
        int boardRow = playerPosition / cells.length;
        int boardColumn = (playerPosition % cells.length);
        return cells[boardRow][boardColumn];
    }
}
