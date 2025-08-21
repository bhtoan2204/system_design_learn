package distributed_system;

import java.io.IOException;

public class Main {
    public static void main(String[] args) {
        try {
            LeaderElection leaderElection = new LeaderElection();
        } catch (IOException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
    }
}