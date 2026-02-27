public class VersionOne {

    private static final Object lock = new Object();
    private static boolean isNumsTurn = false;

    public static void main(String[] args) {
        Thread numThread = new Thread(() -> {
            for (int i = 1; i <= 26; i++) {
                try {
                    synchronized (lock) {
                        while (!isNumsTurn) {
                            // System.out.println("thread A wait()");
                            lock.wait();
                        }
                        System.out.println(i);
                        isNumsTurn = false;
                        // System.out.println("thread A notify()");
                        lock.notify();
                    }
                } catch (Exception e) {
                    Thread.currentThread().interrupt();
                }
            }
        });

        Thread letterThread = new Thread(() -> {
            for (int i = 0; i < 26; i++) {
                try {
                    synchronized (lock) {
                        while (isNumsTurn) {
                            // System.out.println("thread B wait()");
                            lock.wait();
                        }
                        char letter = (char) ('A' + i);
                        System.out.println(letter);
                        isNumsTurn = true;
                        // System.out.println("thread B notify()");
                        lock.notify();
                    }
                } catch (Exception ignored) {
                    Thread.currentThread().interrupt();
                }
            }
        });

        System.out.println("hello world");
        numThread.start();
        letterThread.start();
    }
}
