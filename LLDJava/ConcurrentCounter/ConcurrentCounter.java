package LLDJava.ConcurrentCounter;

import java.util.concurrent.ThreadLocalRandom;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.atomic.LongAccumulator;
import java.util.concurrent.atomic.LongAdder;

public class ConcurrentCounter {

    // Approach 1: AtomicInteger - Good for low to medium contention
    public static class AtomicCounter {
        private final AtomicInteger counter = new AtomicInteger(0);

        public int increment() { return counter.incrementAndGet(); }
        public int decrement() { return counter.decrementAndGet(); }
        public int get() { return counter.get(); }
        public int add(int delta) { return counter.addAndGet(delta); }
        public void set(int value) { counter.set(value); }
    }

    // Approach 2: LongAdder - Excellent for high contention
    public static class LongAdderCounter {
        private final LongAdder counter = new LongAdder();

        public void increment() { counter.increment(); }
        public void decrement() { counter.decrement(); }
        public long get() { return counter.sum(); }
        public void add(long delta) { counter.add(delta); }
        public void reset() { counter.reset(); }
    }

    // Approach 3: Custom sharded counter for extreme contention
    public static class ShardedCounter {
        private static final int SHARD_COUNT = 16; // Must be a power of 2
        private final AtomicLong[] shards;

        public ShardedCounter() {
            this.shards = new AtomicLong[SHARD_COUNT];
            for (int i = 0; i < SHARD_COUNT; i++) {
                shards[i] = new AtomicLong(0);
            }
        }

        public void increment() {
            int shardIndex = getShardIndex();
            shards[shardIndex].incrementAndGet();
        }

        public void decrement() {
            int shardIndex = getShardIndex();
            shards[shardIndex].decrementAndGet();
        }

        public long get() {
            long total = 0;
            for (AtomicLong shard : shards) {
                total += shard.get();
            }
            return total;
        }

        public void add(long delta) {
            int shardIndex = getShardIndex();
            shards[shardIndex].addAndGet(delta);
        }

        private int getShardIndex() {
            int threadHash = System.identityHashCode(Thread.currentThread());
            int rnd = ThreadLocalRandom.current().nextInt();
            return (threadHash ^ rnd) & (SHARD_COUNT - 1);
        }
    }

    // Approach 4: LongAccumulator for custom ops
    public static class AccumulatorCounter {
        private final LongAccumulator counter = new LongAccumulator(Long::sum, 0);

        public void increment() { counter.accumulate(1); }
        public void decrement() { counter.accumulate(-1); }
        public long get() { return counter.longValue(); }
        public void add(long delta) { counter.accumulate(delta); }
        public void reset() { counter.reset(); }
    }

    // ---------------- Benchmarking Utility ----------------

    public static void benchmark(int threadCount, int operationsPerThread) throws InterruptedException {
        System.out.println("Benchmarking with " + threadCount + " threads, " + operationsPerThread + " operations each");

        // AtomicCounter
        AtomicCounter atomicCounter = new AtomicCounter();
        long startTime = System.currentTimeMillis();
        runBenchmark(atomicCounter, threadCount, operationsPerThread);
        long atomicTime = System.currentTimeMillis() - startTime;

        // LongAdder
        LongAdderCounter longAdderCounter = new LongAdderCounter();
        startTime = System.currentTimeMillis();
        runBenchmark(longAdderCounter, threadCount, operationsPerThread);
        long longAdderTime = System.currentTimeMillis() - startTime;

        // ShardedCounter
        ShardedCounter shardedCounter = new ShardedCounter();
        startTime = System.currentTimeMillis();
        runBenchmark(shardedCounter, threadCount, operationsPerThread);
        long shardedTime = System.currentTimeMillis() - startTime;

        System.out.println("AtomicInteger: " + atomicTime + "ms");
        System.out.println("LongAdder: " + longAdderTime + "ms");
        System.out.println("ShardedCounter: " + shardedTime + "ms");

        System.out.println("Final values:");
        System.out.println("AtomicInteger: " + atomicCounter.get());
        System.out.println("LongAdder: " + longAdderCounter.get());
        System.out.println("ShardedCounter: " + shardedCounter.get());
    }

    public static void runBenchmark(Object counter, int threadCount, int operationsPerThread) throws InterruptedException {
        Thread[] threads = new Thread[threadCount];

        for (int i = 0; i < threadCount; i++) {
            threads[i] = new Thread(() -> {
                for (int j = 0; j < operationsPerThread; j++) {
                    if (counter instanceof AtomicCounter) {
                        ((AtomicCounter) counter).increment();
                    } else if (counter instanceof LongAdderCounter) {
                        ((LongAdderCounter) counter).increment();
                    } else if (counter instanceof ShardedCounter) {
                        ((ShardedCounter) counter).increment();
                    }
                }
            });
        }

        for (Thread t : threads) t.start();
        for (Thread t : threads) t.join();
    }

    // ---------------- Main Method ----------------

    public static void main(String[] args) throws InterruptedException {
        System.out.println("=== Concurrent Counter Benchmarking ===\n");

        int[][] scenarios = {
                {10, 100000},
                {25, 40000},
                {50, 20000},
                {100, 10000},
                {200, 5000},
                {500, 2000},
                {1000, 1000}
        };

        for (int[] scenario : scenarios) {
            int threadCount = scenario[0];
            int operationsPerThread = scenario[1];

            System.out.println(">>>>>>>>>> Testing: " + threadCount +
                    " threads, " + operationsPerThread + " ops/thread <<<<<<<<<<");

            benchmark(threadCount, operationsPerThread);
            System.out.println("---\n");
            Thread.sleep(100);
        }
    }
}
