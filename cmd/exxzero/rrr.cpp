// #include <iostream>
// #include <vector>
// #include <algorithm>

// using namespace std;

// bool canFillWater(const vector<long long>& heights, int n, long long k, long long h) {
//     long long freeSquares = 0;
//     for (int i = 0; i < n; i++) {
//         if (heights[i] < h) {
//             freeSquares += h - heights[i];
//         }
//         if (freeSquares >= k) return true;
//     }
//     return false;
// }

// int main() {

//     int n;
//     long long k;
//     cin >> n >> k;
//     vector<long long> heights(n);

//     for (int i = 0; i < n; i++) {
//         cin >> heights[i];
//     }

//     long long low = 0, high = 2e18; 
//     long long result = high;

//     while (high - low > 1) {
//         long long mid = low + (high - low) / 2;
//         if (canFillWater(heights, n, k, mid)) {
//             result = mid;
//             high = mid;
//         }
//         else {
//             low = mid;
//         }
//     }

//     cout << result << "\n";
//     return 0;
// }
#include <iostream>
#include <vector>
#include <algorithm>

using namespace std;

bool canFillWater(const vector<long long>& heights, int n, long long k, long long h) {
    long long freeSquares = 0;
    for (int i = 0; i < n; i++) {
        if (heights[i] < h) {
            freeSquares += h - heights[i];
        }
        if (freeSquares >= k) return true;
    }
    return false;
}

long long findMinHeight(const vector<long long>& heights, int n, long long k) {
    long long low = 0, high = 2e18;
    while (high - low > 1) {
        long long mid = low + (high - low) / 2;
        if (canFillWater(heights, n, k, mid)) {
            high = mid;
        } else {
            low = mid;
        }
    }
    return high;
}

int main() {
    int n;
    long long k;
    cin >> n >> k;
    vector<long long> heights(n);

    for (int i = 0; i < n; i++) {
        cin >> heights[i];
    }

    cout << findMinHeight(heights, n, k) << "\n";
    return 0;
}

