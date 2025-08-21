package distributed_system;

import java.io.IOException;

import org.apache.zookeeper.Watcher;
import org.apache.zookeeper.ZooKeeper;
import org.apache.zookeeper.WatchedEvent;

public class LeaderElection implements Watcher {
    private static final String ZOOKEEPER_ADDRESS = "localhost:2181";
    private static final int SESSION_TIMEOUT = 3000;
    private static final String ELECTION_NAMESPACE = "/election";

    private ZooKeeper zooKeeper;
    private String currentZnodeName;

    public LeaderElection() throws IOException {
        this.zooKeeper = new ZooKeeper(ZOOKEEPER_ADDRESS, SESSION_TIMEOUT, this);
    }

    @Override
    public void process(WatchedEvent event) {
        switch (event.getType()) {
            case None:
                if (event.getState() == Event.KeeperState.SyncConnected) {
                    System.out.println("Successfully connected to ZooKeeper");
                }
                break;
            default:
                break;
        }
    }
}
