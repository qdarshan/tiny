import java.util.concurrent.Semaphore;

public class VersionTwo {

    public static void main(String[] args) {
        Semaphore numSemaphore = new Semaphore(1);
        Semaphore letterSemaphore = new Semaphore(0);

        Thread numThread = new Thread(() -> {
            for (int i = 1; i <= 26; i++) {
                try {
                    numSemaphore.acquire();
                    System.out.println(i);
                    letterSemaphore.release();
                } catch (Exception ignored) {}
            }
        });

        Thread letterThread = new Thread(() -> {
            for (int i = 1; i <= 26; i++) {
                try {
                    letterSemaphore.acquire();
                    char letter = (char) ('A' + i - 1);
                    System.out.println(letter);
                    numSemaphore.release();
                } catch (Exception ignored) {}
            }
        });

        numThread.start();
        letterThread.start();
    }
}
