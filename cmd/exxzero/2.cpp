#include <iostream>
#include <vector>
#include <limits>
#include <set>
#include <algorithm>

using namespace std;

const double INF = numeric_limits<double>::infinity();

void dfs(int start, int current, const vector<vector<double>>& weights, set<int>& visited, vector<int>& path, double cost) {
    int n = weights.size();
    
    if (visited.size() == n) {
        if (weights[current][start] != INF) {
            path.push_back(start + 1);
            cost += weights[current][start];
            
            for (int i = 0; i < path.size(); ++i) {
                if (i > 0) cout << " → ";
                cout << path[i];
            }
            
            cout << "\n";
            for (int i = 0; i < path.size() - 1; ++i) {
                if (i > 0) cout << " + ";
                cout << weights[path[i] - 1][path[i + 1] - 1];
            }
            //cout << " + " << weights[current][start];
            cout << " = " << cost << "\n\n";
            
            path.pop_back();
        }
        return;
    }
    
    vector<pair<int, double>> min_edges;
    double min_cost = INF;
    
    for (int neighbor = 0; neighbor < n; ++neighbor) {
        if (visited.find(neighbor) == visited.end() && weights[current][neighbor] < min_cost) {
            min_edges.clear();
            min_edges.emplace_back(neighbor, weights[current][neighbor]);
            min_cost = weights[current][neighbor];
        } else if (visited.find(neighbor) == visited.end() && weights[current][neighbor] == min_cost) {
            min_edges.emplace_back(neighbor, weights[current][neighbor]);
        }
    }
    
    for (const auto& [neighbor, edge_cost] : min_edges) {
        visited.insert(neighbor);
        path.push_back(neighbor + 1);
        dfs(start, neighbor, weights, visited, path, cost + edge_cost);
        visited.erase(neighbor);
        path.pop_back();
    }
}

void find_cycles(const vector<vector<double>>& weights) {
    int n = weights.size();
    
    for (int start = 0; start < n; ++start) {
        cout << "Початкова вершина " << start + 1 << ":\n";
        set<int> visited = { start };
        vector<int> path = { start + 1 };
        dfs(start, start, weights, visited, path, 0.0);
    }
}

int main() {
    vector<vector<double>> weights = {
        { INF, 4, 6, 5, 5, 5, 4, 5 },
        { 4, INF, 7, 1, 6, 5, 1, 4 },
        { 6, 7, INF, 1, 1, 1, 5, 7 },
        { 5, 1, 1, INF, 6, 5, 1, 7 },
        { 5, 6, 1, 6, INF, 5, 5, 6 },
        { 5, 5, 1, 5, 5, INF, 5, 7 },
        { 4, 1, 5, 1, 5, 5, INF, 5 },
        { 5, 4, 7, 7, 6, 7, 5, INF }
    };
    
    find_cycles(weights);
    
    return 0;
}