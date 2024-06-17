#include <iostream>
#include <vector>
#include <limits>
#include <set>
#include <algorithm>

using namespace std;

const double INF = numeric_limits<double>::infinity();

void dfs(int start, int current, const vector<vector<double>>& weights, set<int>& visited, vector<int>& path, double cost, double& min_cost, vector<int>& min_path) {
    int n = weights.size();
    
    if (visited.size() == n) {
        if (weights[current][start] != INF) {
            path.push_back(start + 1);
            cost += weights[current][start];
            
            if (cost < min_cost) {
                min_cost = cost;
                min_path = path;
            }
            
            path.pop_back();
        }
        return;
    }
    
    vector<pair<int, double>> min_edges;
    double min_edge_cost = INF;
    
    for (int neighbor = 0; neighbor < n; ++neighbor) {
        if (visited.find(neighbor) == visited.end() && weights[current][neighbor] < min_edge_cost) {
            min_edges.clear();
            min_edges.emplace_back(neighbor, weights[current][neighbor]);
            min_edge_cost = weights[current][neighbor];
        } else if (visited.find(neighbor) == visited.end() && weights[current][neighbor] == min_edge_cost) {
            min_edges.emplace_back(neighbor, weights[current][neighbor]);
        }
    }
    
    for (const auto& [neighbor, edge_cost] : min_edges) {
        visited.insert(neighbor);
        path.push_back(neighbor + 1);
        dfs(start, neighbor, weights, visited, path, cost + edge_cost, min_cost, min_path);
        visited.erase(neighbor);
        path.pop_back();
    }
}

void find_cycles(const vector<vector<double>>& weights) {
    int n = weights.size();
    
    for (int start = 0; start < n; ++start) {
        cout << "Починаючи з вершини " << start + 1 << ":\n";
        set<int> visited = { start };
        vector<int> path = { start + 1 };
        double min_cost = INF;
        vector<int> min_path;
        dfs(start, start, weights, visited, path, 0.0, min_cost, min_path);
        
        if (!min_path.empty()) {
            cout << "Мінімальний шлях: ";
            for (int i = 0; i < min_path.size(); ++i) {
                if (i > 0) cout << "→";
                cout << min_path[i];
            }
            cout << "\n";
            
            cout << "Ваги: ";
            for (int i = 0; i < min_path.size() - 1; ++i) {
                if (i > 0) cout << " + ";
                cout << weights[min_path[i] - 1][min_path[i + 1] - 1];
            }
            cout << " = " << min_cost << "\n\n";
        }
    }
}

int main() {
    setlocale(LC_ALL, "");
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
