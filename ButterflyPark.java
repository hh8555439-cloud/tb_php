import java.util.*;

public class ButterflyPark {
    public static void main(String[] args) {
        Scanner in = new Scanner(System.in);
        int T = in.nextInt();
        StringBuilder sb = new StringBuilder();

        while (T-- > 0) {
            int n = in.nextInt();
            long k = in.nextLong();

            long[] prefix = new long[n + 1];
            for (int i = 1; i <= n; i++) {
                prefix[i] = prefix[i - 1] + in.nextInt();
            }

            long[] S = new long[n + 1];
            for (int i = 1; i <= n; i++) {
                S[i] = S[i - 1] + prefix[i];
            }

            long lo = 0, hi = S[n];
            while (lo < hi) {
                long mid = lo + (hi - lo) / 2;
                if (countLE(n, mid, prefix, S) >= k) {
                    hi = mid;
                } else {
                    lo = mid + 1;
                }
            }

            sb.append(lo).append('\n');
        }

        System.out.print(sb);
    }

    static long countLE(int n, long v, long[] prefix, long[] S) {
        long cnt = 0;
        for (int l = 1; l <= n; l++) {
            long C = S[l - 1] - (long)(l - 1) * prefix[l - 1];
            long threshold = v + C;
            int low = l, high = n, ans = l - 1;
            while (low <= high) {
                int mid = (low + high) >>> 1;
                if (S[mid] - (long)mid * prefix[l - 1] <= threshold) {
                    ans = mid;
                    low = mid + 1;
                } else {
                    high = mid - 1;
                }
            }
            cnt += ans - l + 1;
        }
        return cnt;
    }
}
