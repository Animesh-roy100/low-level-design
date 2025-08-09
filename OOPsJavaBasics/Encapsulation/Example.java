package OOPsJavaBasics.Encapsulation;

class CricketScorer {
    private int score = 0;

    public int getScore() {
        return score;
    }

    public void addScore(int score) {
        this.score += score;
    }

    public void four() {
        addScore(4);
    }

    public void single() {
        addScore(1);
    }

    public void six() {
        addScore(6);
    }

    public void printScore() {
        System.out.println("Score: " + score);
    }
}

public class Example {
    public static void main(String[] args) {
        CricketScorer cricketScorer = new CricketScorer();

        cricketScorer.four();
		cricketScorer.four();
		cricketScorer.single();
		cricketScorer.six();
		cricketScorer.six();
		cricketScorer.six();
		cricketScorer.printScore();
    }
}
