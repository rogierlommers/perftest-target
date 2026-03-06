#!/bin/bash

# ============================================================
# Vegeta Parallel Load Test Script
# Requirements: vegeta (https://github.com/tsenart/vegeta)
#   macOS:  brew install vegeta
#   Linux:  https://github.com/tsenart/vegeta/releases
# ============================================================

set -euo pipefail

DURATION="30s"
OUTPUT_DIR="./vegeta-results"
BASE="http://localhost:3000"

mkdir -p "$OUTPUT_DIR"
rm -f "$OUTPUT_DIR"/*.bin "$OUTPUT_DIR"/*.txt

command -v vegeta >/dev/null 2>&1 || { echo "vegeta is not installed. Install with: brew install vegeta"; exit 1; }

# ---- Helper: run a single attack in the background ----
attack() {
  local method="$1"
  local path="$2"
  local rate="$3"
  local name="${method} ${path}"
  local safe_name="${method}_${path#/}"

  echo "▶ Starting attack: $name @ ${rate} req/s for $DURATION"

  # Write a target file with a single entry
  local target_file
  target_file=$(mktemp)
  printf '%s %s%s\n' "$method" "$BASE" "$path" > "$target_file"

  (
    vegeta attack \
      -targets="$target_file" \
      -rate="$rate" \
      -duration="$DURATION" \
      -name="$name" \
      > "$OUTPUT_DIR/${safe_name}.bin"

    rm -f "$target_file"

    vegeta report -type=text < "$OUTPUT_DIR/${safe_name}.bin" > "$OUTPUT_DIR/${safe_name}.report.txt"
    vegeta report -type='hist[0,5ms,10ms,25ms,50ms,100ms,250ms,500ms]' \
      < "$OUTPUT_DIR/${safe_name}.bin" > "$OUTPUT_DIR/${safe_name}.hist.txt"

    echo "✔ Done: $name"
  ) &
}

# ============================================================
# Define endpoints with random rates (10–1100 req/s)
# ============================================================

attack "GET"  "/users"     $(( (RANDOM % 1091) + 10 ))
attack "POST" "/users"     $(( (RANDOM % 1091) + 10 ))
attack "GET"  "/tasks"     $(( (RANDOM % 1091) + 10 ))
attack "GET"  "/documents" $(( (RANDOM % 1091) + 10 ))
attack "POST" "/documents" $(( (RANDOM % 1091) + 10 ))

# ============================================================
# Wait for all background attacks to complete
# ============================================================
wait
echo ""
echo "========================================"
echo "       LOAD TEST RESULTS SUMMARY"
echo "========================================"

for bin_file in "$OUTPUT_DIR"/*.bin; do
  name=$(basename "$bin_file" .bin | tr '_' ' ')
  report_file="${bin_file%.bin}.report.txt"
  echo ""
  echo "--- $name ---"
  cat "$report_file"
done

echo ""
echo "Full results saved to: $OUTPUT_DIR/"
echo "  *.bin         — raw binary (pipe to 'vegeta report' or 'vegeta plot')"
echo "  *.report.txt  — text summary"
echo "  *.hist.txt    — latency histogram"
echo ""
echo "Tip: generate an HTML plot with:"
echo "  vegeta plot $OUTPUT_DIR/*.bin > results.html && open results.html"