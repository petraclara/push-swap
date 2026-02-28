/* ─────────────────────────────────────────────
   push_swap visualizer — app.js
   Reads ops.txt from the same directory (no prompts, no popups).
   Fixes:
     - no layout-shift / scroll shaking during animation
     - ops.txt read directly via fetch
     - input taken from text field in the UI
───────────────────────────────────────────── */

"use strict";

// ── OP METADATA ──────────────────────────────────────────────────
const OP_DESC = {
  sa: "swap top 2 of A",
  sb: "swap top 2 of B",
  ss: "sa + sb",
  pa: "push B → A",
  pb: "push A → B",
  ra: "rotate A",
  rb: "rotate B",
  rr: "ra + rb",
  rra: "rev-rotate A",
  rrb: "rev-rotate B",
  rrr: "rra + rrb",
};

// ── STACK SIMULATION ─────────────────────────────────────────────
function applyOp(stateA, stateB, op) {
  const a = [...stateA],
    b = [...stateB];
  switch (op) {
    case "sa":
      if (a.length >= 2) [a[0], a[1]] = [a[1], a[0]];
      break;
    case "sb":
      if (b.length >= 2) [b[0], b[1]] = [b[1], b[0]];
      break;
    case "ss":
      if (a.length >= 2) [a[0], a[1]] = [a[1], a[0]];
      if (b.length >= 2) [b[0], b[1]] = [b[1], b[0]];
      break;
    case "pa":
      if (b.length) a.unshift(b.shift());
      break;
    case "pb":
      if (a.length) b.unshift(a.shift());
      break;
    case "ra":
      if (a.length >= 2) a.push(a.shift());
      break;
    case "rb":
      if (b.length >= 2) b.push(b.shift());
      break;
    case "rr":
      if (a.length >= 2) a.push(a.shift());
      if (b.length >= 2) b.push(b.shift());
      break;
    case "rra":
      if (a.length >= 2) a.unshift(a.pop());
      break;
    case "rrb":
      if (b.length >= 2) b.unshift(b.pop());
      break;
    case "rrr":
      if (a.length >= 2) a.unshift(a.pop());
      if (b.length >= 2) b.unshift(b.pop());
      break;
  }
  return [a, b];
}

function isSorted(arr) {
  for (let i = 0; i < arr.length - 1; i++)
    if (arr[i] > arr[i + 1]) return false;
  return true;
}

// ── STATE ────────────────────────────────────────────────────────
let allOps = []; // string[]
let snapshots = []; // Array of { a: number[], b: number[] }
let stepIdx = 0; // current snapshot index
let totalN = 0; // initial stack size (for color mapping)
let playing = false;
let playTimer = null;

// ── DOM REFS ─────────────────────────────────────────────────────
const dom = {
  input: document.getElementById("stack-input"),
  btnLoad: document.getElementById("btn-load"),
  btnPlay: document.getElementById("btn-play"),
  btnPrev: document.getElementById("btn-prev"),
  btnNext: document.getElementById("btn-next"),
  btnReset: document.getElementById("btn-reset"),
  speedSlider: document.getElementById("speed"),
  speedVal: document.getElementById("speed-val"),
  stepCur: document.getElementById("step-current"),
  stepTot: document.getElementById("step-total"),
  sortedBadge: document.getElementById("sorted-badge"),
  barsA: document.getElementById("bars-a"),
  barsB: document.getElementById("bars-b"),
  countA: document.getElementById("count-a"),
  countB: document.getElementById("count-b"),
  opPill: document.getElementById("op-pill"),
  opDesc: document.getElementById("op-desc"),
  logBody: document.getElementById("log-body"),
  totalOpsLbl: document.getElementById("total-ops-label"),
  stA: document.getElementById("st-a"),
  stB: document.getElementById("st-b"),
  stOp: document.getElementById("st-op"),
  stSorted: document.getElementById("st-sorted"),
};

// ── COLORS ───────────────────────────────────────────────────────
function barColor(indexedVal, total, isA) {
  if (total <= 1) return isA ? "#00d4ff" : "#ff6b35";
  const t = indexedVal / Math.max(1, total - 1);
  if (isA) {
    const r = Math.round(t * 40);
    const g = Math.round(212 - t * 130);
    const b = Math.round(255 - t * 80);
    return `rgb(${r},${g},${b})`;
  } else {
    const r = 255;
    const g = Math.round(107 - t * 80);
    const b = Math.round(53 - t * 40);
    return `rgb(${r},${g},${b})`;
  }
}

// ── RENDER STACKS ────────────────────────────────────────────────
// We pre-calculate bar height once and never change it → no reflow jitter.
let barHeight = 24;
let barsAreaH = 0;

function calcBarHeight(maxCount) {
  const el = dom.barsA;
  const h = el.clientHeight || 400;
  barsAreaH = h;
  const gap = 2;
  const rows = Math.max(maxCount, 1);
  barHeight = Math.max(14, Math.min(32, Math.floor((h - gap * rows) / rows)));
}

function renderStacks(snap, prevSnap) {
  renderOneSide(dom.barsA, snap.a, prevSnap ? prevSnap.a : null, true);
  renderOneSide(dom.barsB, snap.b, prevSnap ? prevSnap.b : null, false);
}

function renderOneSide(container, arr, prevArr, isA) {
  // Diff: only re-render if content changed
  const key = arr.join(",");
  if (container._lastKey === key) return;
  container._lastKey = key;

  // Remove old bars (keep .empty-msg in DOM)
  const oldBars = container.querySelectorAll(".bar-row");
  oldBars.forEach((el) => el.remove());

  const emptyMsg = container.querySelector(".empty-msg");
  if (!arr.length) {
    if (emptyMsg) emptyMsg.style.display = "flex";
    return;
  }
  if (emptyMsg) emptyMsg.style.display = "none";

  const vpW = container.clientWidth - 60; // leave room for index labels
  const frag = document.createDocumentFragment();

  // arr[0] = top of stack → first DOM child → visually at top.
  // flex-direction: column with justify-content: flex-start.
  for (let i = 0; i < arr.length; i++) {
    const val = arr[i];
    const row = document.createElement("div");
    row.className = "bar-row" + (i === 0 ? " is-top" : "");

    const idxLbl = document.createElement("span");
    idxLbl.className = "bar-idx";
    idxLbl.textContent = i === 0 ? "▶" : i;
    row.appendChild(idxLbl);

    const bar = document.createElement("div");
    bar.className = "bar";
    const w = Math.max(24, Math.round(((val + 1) / totalN) * vpW));
    bar.style.width = w + "px";
    bar.style.height = barHeight + "px";
    bar.style.backgroundColor = barColor(val, totalN, isA);

    const lbl = document.createElement("span");
    lbl.className = "bar-val";
    lbl.textContent = val;
    bar.appendChild(lbl);
    row.appendChild(bar);
    frag.appendChild(row);
  }

  container.appendChild(frag);
  // arr[0] is at the top of the DOM — scroll there so the active element is visible
  container.scrollTop = 0;
}

// ── RENDER STEP ──────────────────────────────────────────────────
function renderStep(idx) {
  idx = Math.max(0, Math.min(idx, snapshots.length - 1));
  stepIdx = idx;

  const snap = snapshots[idx];
  const prevSnap = idx > 0 ? snapshots[idx - 1] : null;
  const opIdx = idx - 1;
  const op = opIdx >= 0 ? allOps[opIdx] : null;

  renderStacks(snap, prevSnap);

  // counts
  dom.countA.textContent = snap.a.length;
  dom.countB.textContent = snap.b.length;
  dom.stA.textContent = snap.a.length;
  dom.stB.textContent = snap.b.length;
  dom.stepCur.textContent = idx;

  // op display
  if (op) {
    dom.opPill.textContent = op;
    dom.opPill.className = "op-pill op-" + op;
    dom.opDesc.textContent = OP_DESC[op] || "";
  } else {
    dom.opPill.textContent = "";
    dom.opPill.className = "op-pill";
    dom.opDesc.textContent = "— initial state —";
  }

  dom.stOp.textContent = op || "—";

  // sorted badge
  const sorted = isSorted(snap.a) && snap.b.length === 0 && snap.a.length > 0;
  dom.sortedBadge.textContent = sorted ? "✓ sorted" : "sorting…";
  dom.sortedBadge.className = "badge" + (sorted ? " sorted" : "");
  dom.stSorted.textContent = sorted ? "✓ sorted" : "not sorted";
  dom.stSorted.className = "st-sorted" + (sorted ? " ok" : "");

  // log highlight (without causing layout shift)
  highlightLogEntry(opIdx);

  updateButtons();
}

// ── LOG ──────────────────────────────────────────────────────────
function buildLog() {
  dom.logBody.innerHTML = "";
  const frag = document.createDocumentFragment();
  allOps.forEach((op, i) => {
    const el = document.createElement("div");
    el.className = "log-entry";
    el.dataset.i = i;

    const num = document.createElement("span");
    num.className = "log-num";
    num.textContent = i + 1;

    const opEl = document.createElement("span");
    opEl.className = "log-op op-" + op;
    opEl.textContent = op;

    const desc = document.createElement("span");
    desc.className = "log-desc";
    desc.textContent = OP_DESC[op] || "";

    el.appendChild(num);
    el.appendChild(opEl);
    el.appendChild(desc);

    // click to jump
    el.addEventListener("click", () => {
      stopPlay();
      renderStep(i + 1);
    });

    frag.appendChild(el);
  });
  dom.logBody.appendChild(frag);
  dom.totalOpsLbl.textContent = allOps.length + " ops";
  dom.stepTot.textContent = allOps.length;
}

let _lastHighlightedIdx = -2;
function highlightLogEntry(opIdx) {
  if (opIdx === _lastHighlightedIdx) return;
  _lastHighlightedIdx = opIdx;

  const entries = dom.logBody.querySelectorAll(".log-entry");
  entries.forEach((el, i) => {
    el.classList.toggle("past", i < opIdx);
    el.classList.toggle("current", i === opIdx);
  });

  // scroll log — but NEVER scroll the main viewport
  if (opIdx >= 0 && entries[opIdx]) {
    entries[opIdx].scrollIntoView({ block: "nearest", behavior: "smooth" });
  }
}

// ── PLAY / PAUSE ─────────────────────────────────────────────────
function startPlay() {
  if (stepIdx >= snapshots.length - 1) renderStep(0);
  playing = true;
  dom.btnPlay.textContent = "⏸ pause";
  dom.btnPlay.classList.add("paused");
  tick();
}

function stopPlay() {
  playing = false;
  clearTimeout(playTimer);
  dom.btnPlay.textContent = "▶ play";
  dom.btnPlay.classList.remove("paused");
}

function tick() {
  if (!playing) return;
  if (stepIdx >= snapshots.length - 1) {
    stopPlay();
    return;
  }
  renderStep(stepIdx + 1);
  const speed = parseInt(dom.speedSlider.value, 10);
  // map speed 1-10 → delay ~1200ms … ~20ms
  const delay = Math.max(16, Math.round(1400 / (speed * speed)));
  playTimer = setTimeout(tick, delay);
}

// ── BUTTONS ──────────────────────────────────────────────────────
function updateButtons() {
  const has = snapshots.length > 1;
  dom.btnPlay.disabled = !has;
  dom.btnPrev.disabled = !has || stepIdx <= 0;
  dom.btnNext.disabled = !has || stepIdx >= snapshots.length - 1;
  dom.btnReset.disabled = !has;
}

dom.btnPlay.addEventListener("click", () => {
  if (!playing) startPlay();
  else stopPlay();
});

dom.btnPrev.addEventListener("click", () => {
  stopPlay();
  renderStep(stepIdx - 1);
});

dom.btnNext.addEventListener("click", () => {
  stopPlay();
  renderStep(stepIdx + 1);
});

dom.btnReset.addEventListener("click", () => {
  stopPlay();
  renderStep(0);
});

dom.speedSlider.addEventListener("input", function () {
  dom.speedVal.textContent = this.value;
});

// ── LOAD ─────────────────────────────────────────────────────────
function initFromOps(opsText, initialNums) {
  const ops = opsText
    .trim()
    .split("\n")
    .map((s) => s.trim())
    .filter(Boolean)
    .filter((s) => s in OP_DESC);

  if (!ops.length) {
    showError("ops.txt is empty or contains no valid operations.");
    return;
  }

  allOps = ops;
  totalN = initialNums.length;

  // simulate all states
  let a = [...initialNums],
    b = [];
  snapshots = [{ a: [...a], b: [] }];
  for (const op of ops) {
    [a, b] = applyOp(a, b, op);
    snapshots.push({ a: [...a], b: [...b] });
  }

  // layout
  calcBarHeight(totalN);
  buildLog();

  // invalidate cached keys
  dom.barsA._lastKey = null;
  dom.barsB._lastKey = null;

  renderStep(0);
}

function parseInput(raw) {
  const nums = raw.trim().split(/\s+/).map(Number);
  if (nums.some(isNaN)) return null;
  const unique = new Set(nums);
  if (unique.size !== nums.length) return null;
  return nums;
}

dom.btnLoad.addEventListener("click", async () => {
  await fetchAndLoad();
});

// ── AUTO-LOAD on page open ────────────────────────────────────────
// Fetch both ops.txt and initial.txt written by push_swap.
// The input field is only used as a fallback if initial.txt is missing.
window.addEventListener("load", async () => {
  await fetchAndLoad();
  updateButtons();
});

async function fetchAndLoad() {
  // 1. Fetch ops.txt
  let opsText = "";
  try {
    const resp = await fetch("ops.txt");
    if (!resp.ok) throw new Error();
    opsText = await resp.text();
    if (!opsText.trim()) throw new Error("empty");
  } catch (_) {
    // Not served via HTTP yet — user will click Load after starting the server
    return;
  }

  // 2. Fetch initial.txt (written by push_swap alongside ops.txt)
  let nums = null;
  try {
    const resp = await fetch("initial.txt");
    if (resp.ok) {
      const text = await resp.text();
      nums = parseInput(text);
    }
  } catch (_) {
    /* fall through to input field */
  }

  // 3. Fallback: read from the text input if initial.txt wasn't found
  if (!nums) {
    nums = parseInput(dom.input.value);
  }

  if (!nums || nums.length < 2) {
    showError(
      "Could not determine initial stack. Enter numbers in the input field and click Load.",
    );
    return;
  }

  // 4. Sync the input field so user can see what was loaded
  dom.input.value = nums.join(" ");

  stopPlay();
  initFromOps(opsText, nums);
}

// ── KEYBOARD ─────────────────────────────────────────────────────
document.addEventListener("keydown", (e) => {
  if (e.target === dom.input) return;
  switch (e.key) {
    case " ":
      e.preventDefault();
      dom.btnPlay.click();
      break;
    case "ArrowRight":
      e.preventDefault();
      dom.btnNext.click();
      break;
    case "ArrowLeft":
      e.preventDefault();
      dom.btnPrev.click();
      break;
    case "r":
    case "R":
      dom.btnReset.click();
      break;
  }
});

// ── RESIZE ───────────────────────────────────────────────────────
// Recalculate bar heights on resize without triggering scroll
let resizeTimer;
window.addEventListener("resize", () => {
  clearTimeout(resizeTimer);
  resizeTimer = setTimeout(() => {
    if (!snapshots.length) return;
    calcBarHeight(totalN);
    // Force re-render by clearing cache
    dom.barsA._lastKey = null;
    dom.barsB._lastKey = null;
    const snap = snapshots[stepIdx];
    renderOneSide(dom.barsA, snap.a, null, true);
    renderOneSide(dom.barsB, snap.b, null, false);
  }, 100);
});

// ── ERROR TOAST ──────────────────────────────────────────────────
function showError(msg) {
  const existing = document.querySelector(".error-toast");
  if (existing) existing.remove();
  const el = document.createElement("div");
  el.className = "error-toast";
  el.textContent = msg;
  document.body.appendChild(el);
  setTimeout(() => el.remove(), 4000);
}
