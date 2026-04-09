import { useState, useEffect, useRef } from 'react'
import { api } from '../api'

const MOOD_EMOJI = ['', '😞','😔','😕','😐','🙂','😊','😄','😁','🤩','🥳']
const MOOD_COLOR = ['','#e05c5c','#e07c5c','#e09c5c','#c9a96e','#b8b06e','#96b86e','#6dbf8a','#5cbfaa','#5caebf','#5c8ebf']

function formatDate(str) {
  const d = new Date(str)
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function avgMood(entries) {
  if (!entries.length) return '—'
  return (entries.reduce((s, e) => s + e.mood_score, 0) / entries.length).toFixed(1)
}

function bestDay(entries) {
  if (!entries.length) return '—'
  const days = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat']
  const map = {}
  entries.forEach(e => {
    const d = new Date(e.created_at).getDay()
    if (!map[d]) map[d] = { sum: 0, count: 0 }
    map[d].sum += e.mood_score
    map[d].count++
  })
  let best = null, bestAvg = 0
  Object.entries(map).forEach(([d, v]) => {
    const a = v.sum / v.count
    if (a > bestAvg) { bestAvg = a; best = d }
  })
  return best !== null ? days[best] : '—'
}

function moodTrend(entries) {
  if (entries.length < 2) return null
  const recent = entries.slice(0, 3)
  const older = entries.slice(3, 6)
  if (!older.length) return null
  const r = recent.reduce((s, e) => s + e.mood_score, 0) / recent.length
  const o = older.reduce((s, e) => s + e.mood_score, 0) / older.length
  const diff = r - o
  if (Math.abs(diff) < 0.3) return { dir: 'stable', text: 'Stable', color: '#c9a96e' }
  return diff > 0
    ? { dir: 'up', text: `+${diff.toFixed(1)} vs last week`, color: '#6dbf8a' }
    : { dir: 'down', text: `${diff.toFixed(1)} vs last week`, color: '#e05c5c' }
}

function InsightPanel({ entries }) {
  const [insight, setInsight] = useState('')
  const [loading, setLoading] = useState(false)
  const [done, setDone] = useState(false)

  const generate = async () => {
    if (entries.length < 3) return
    setLoading(true)
    setInsight('')
    setDone(false)
    try {
      const data = await api.getInsight()
      setInsight(data.insight || 'No insight available.')
      setDone(true)
    } catch {
      setInsight('Could not generate insight. Make sure you have entries with sentiment analyzed.')
      setDone(true)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="insight-panel">
      <div className="insight-header">
        <div className="insight-title">
          <span className="insight-icon">✦</span>
          AI Weekly Insight
        </div>
        <button
          className="insight-btn"
          onClick={generate}
          disabled={loading || entries.length < 3}
          title={entries.length < 3 ? 'Need at least 3 entries' : ''}
        >
          {loading ? <span className="insight-spinner" /> : done ? 'Regenerate' : 'Generate'}
        </button>
      </div>
      {loading && (
        <div className="insight-loading">
          <span className="insight-spinner" /> Analyzing your entries...
        </div>
      )}
      {insight && !loading && (
        <div className="insight-text">{insight}</div>
      )}
      {!insight && !loading && (
        <div className="insight-placeholder">
          {entries.length < 3
            ? 'Write at least 3 entries to unlock AI insights.'
            : 'Click Generate to analyze your mood patterns with AI.'}
        </div>
      )}
    </div>
  )
}

export default function Dashboard({ onLogout }) {
  const [entries, setEntries] = useState([])
  const [loading, setLoading] = useState(true)
  const [content, setContent] = useState('')
  const [mood, setMood] = useState(5)
  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState('')
  const [deletingId, setDeletingId] = useState(null)
  const textareaRef = useRef(null)

  const load = async () => {
    try {
      const data = await api.getEntries()
      setEntries(data.ok || [])
    } catch {
      setEntries([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [])

  const submit = async (e) => {
    e.preventDefault()
    if (!content.trim()) return
    setSubmitting(true)
    setFormError('')
    try {
      await api.createEntry({ content: content.trim(), mood_score: mood })
      setContent('')
      setMood(5)
      await load()
    } catch (err) {
      setFormError(err.message)
    } finally {
      setSubmitting(false)
    }
  }

  const remove = async (id) => {
    setDeletingId(id)
    try {
      await api.deleteEntry(id)
      setEntries(prev => prev.filter(e => e.id !== id))
    } catch {}
    finally { setDeletingId(null) }
  }

  const trend = moodTrend(entries)
  const today = new Date().toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })
  const hour = new Date().getHours()
  const greeting = hour < 12 ? 'Good morning.' : hour < 17 ? 'Good afternoon.' : 'Good evening.'

  return (
    <div className="dashboard">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Fraunces:ital,wght@0,300;0,400;1,300;1,400&family=DM+Sans:wght@300;400;500&display=swap');
        :root {
          --bg: #0b0b0c; --bg2: #131315; --bg3: #1a1a1d;
          --border: #242428; --border2: #2e2e33;
          --text: #e8e4dc; --muted: #6b6760; --muted2: #8a8680;
          --accent: #c9a96e; --accent2: #8b7355; --accent3: #e8c990;
          --danger: #e05c5c; --success: #6dbf8a;
          --serif: 'Fraunces', serif; --sans: 'DM Sans', sans-serif;
          --r: 10px; --r-sm: 6px;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { background: var(--bg); color: var(--text); font-family: var(--sans); font-weight: 300; }
        .dashboard { display: grid; grid-template-columns: 220px 1fr; min-height: 100vh; }

        .sidebar {
          background: var(--bg2); border-right: 1px solid var(--border);
          padding: 28px 16px; display: flex; flex-direction: column;
          position: sticky; top: 0; height: 100vh;
        }
        .sidebar-logo { font-family: var(--serif); font-size: 20px; color: var(--accent); padding: 0 10px; margin-bottom: 28px; letter-spacing: -0.3px; }
        .sidebar-link { display: flex; align-items: center; gap: 9px; padding: 9px 10px; border-radius: var(--r-sm); font-size: 13px; color: var(--muted2); cursor: pointer; transition: all 0.15s; }
        .sidebar-link:hover { background: var(--bg3); color: var(--text); }
        .sidebar-link.active { background: var(--bg3); color: var(--text); border: 1px solid var(--border2); }
        .sidebar-bottom { margin-top: auto; padding-top: 16px; border-top: 1px solid var(--border); display: flex; flex-direction: column; gap: 4px; }
        .btn-logout { padding: 9px 10px; border-radius: var(--r-sm); font-size: 13px; font-family: var(--sans); color: var(--muted); background: none; border: none; cursor: pointer; text-align: left; transition: all 0.15s; }
        .btn-logout:hover { background: rgba(224,92,92,0.08); color: var(--danger); }

        .main { padding: 44px 52px; max-width: 860px; }

        .page-header { margin-bottom: 36px; }
        .page-header h1 { font-family: var(--serif); font-size: 44px; font-weight: 300; line-height: 1; color: var(--text); font-style: italic; }
        .page-header p { color: var(--muted); font-size: 13px; margin-top: 8px; }

        .stats-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 32px; }
        .stat-card { background: var(--bg2); border: 1px solid var(--border); border-radius: var(--r); padding: 18px 20px; }
        .stat-label { font-size: 10px; letter-spacing: 0.1em; text-transform: uppercase; color: var(--muted); margin-bottom: 8px; }
        .stat-value { font-family: var(--serif); font-size: 32px; font-weight: 300; color: var(--accent); line-height: 1; }
        .stat-sub { font-size: 11px; color: var(--muted); margin-top: 5px; }
        .stat-trend { font-size: 11px; margin-top: 5px; font-weight: 500; }

        .insight-panel { background: var(--bg2); border: 1px solid var(--border2); border-radius: var(--r); padding: 22px 24px; margin-bottom: 32px; }
        .insight-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
        .insight-title { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 500; color: var(--text); }
        .insight-icon { color: var(--accent); font-size: 14px; }
        .insight-btn { background: var(--bg3); border: 1px solid var(--border2); color: var(--accent); font-family: var(--sans); font-size: 12px; font-weight: 500; padding: 6px 14px; border-radius: var(--r-sm); cursor: pointer; transition: all 0.15s; display: flex; align-items: center; gap: 6px; }
        .insight-btn:hover:not(:disabled) { background: rgba(201,169,110,0.1); border-color: var(--accent2); }
        .insight-btn:disabled { opacity: 0.4; cursor: not-allowed; }
        .insight-spinner { width: 12px; height: 12px; border: 1.5px solid var(--border2); border-top-color: var(--accent); border-radius: 50%; animation: spin 0.7s linear infinite; display: inline-block; }
        @keyframes spin { to { transform: rotate(360deg); } }
        .insight-text { font-size: 14px; font-weight: 300; color: var(--text); line-height: 1.7; font-style: italic; padding-left: 16px; border-left: 2px solid var(--accent2); }
        .insight-placeholder { font-size: 13px; color: var(--muted); font-style: italic; }
        .insight-loading { display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--muted); }

        .new-entry-card { background: var(--bg2); border: 1px solid var(--border); border-radius: var(--r); padding: 24px; margin-bottom: 32px; }
        .new-entry-card h2 { font-family: var(--serif); font-size: 18px; font-weight: 300; margin-bottom: 16px; color: var(--text); }
        .entry-textarea { width: 100%; background: var(--bg3); border: 1px solid var(--border); border-radius: var(--r-sm); padding: 13px 15px; color: var(--text); font-size: 14px; font-weight: 300; line-height: 1.65; outline: none; resize: none; font-family: var(--sans); transition: border-color 0.2s; }
        .entry-textarea:focus { border-color: var(--accent2); }
        .entry-textarea::placeholder { color: var(--muted); }
        .mood-row { display: flex; align-items: center; gap: 16px; margin-top: 14px; }
        .mood-label { font-size: 11px; letter-spacing: 0.08em; text-transform: uppercase; color: var(--muted); white-space: nowrap; }
        .mood-slider { flex: 1; -webkit-appearance: none; height: 3px; border-radius: 2px; outline: none; cursor: pointer; }
        .mood-slider::-webkit-slider-thumb { -webkit-appearance: none; width: 16px; height: 16px; border-radius: 50%; cursor: pointer; transition: transform 0.15s; border: 2px solid var(--bg2); }
        .mood-slider::-webkit-slider-thumb:hover { transform: scale(1.25); }
        .mood-value { font-family: var(--serif); font-size: 20px; color: var(--accent); min-width: 24px; text-align: right; }
        .entry-actions { display: flex; justify-content: space-between; align-items: center; margin-top: 14px; }
        .char-count { font-size: 11px; color: var(--muted); }
        .btn-submit { background: var(--accent); color: #1a1408; font-size: 12px; font-weight: 500; letter-spacing: 0.05em; padding: 9px 22px; border-radius: var(--r-sm); border: none; cursor: pointer; font-family: var(--sans); transition: opacity 0.15s; }
        .btn-submit:hover { opacity: 0.85; }
        .btn-submit:disabled { opacity: 0.35; cursor: not-allowed; }
        .error-msg { background: rgba(224,92,92,0.1); border: 1px solid rgba(224,92,92,0.25); color: #e08080; font-size: 13px; padding: 9px 13px; border-radius: var(--r-sm); margin-bottom: 14px; }

        .entries-header { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 16px; }
        .entries-header h2 { font-family: var(--serif); font-size: 20px; font-weight: 300; color: var(--text); }
        .entries-count { font-size: 12px; color: var(--muted); }

        .entry-item { background: var(--bg2); border: 1px solid var(--border); border-radius: var(--r); padding: 18px 20px; margin-bottom: 10px; display: grid; grid-template-columns: 44px 1fr auto; gap: 16px; align-items: start; transition: border-color 0.2s, opacity 0.2s; animation: fadeUp 0.25s ease; }
        @keyframes fadeUp { from { opacity: 0; transform: translateY(6px); } to { opacity: 1; transform: translateY(0); } }
        .entry-item:hover { border-color: var(--border2); }
        .entry-item.deleting { opacity: 0.4; pointer-events: none; }

        .entry-mood { display: flex; flex-direction: column; align-items: center; gap: 4px; padding-top: 2px; }
        .entry-mood-num { font-family: var(--serif); font-size: 22px; font-weight: 300; line-height: 1; }
        .entry-mood-emoji { font-size: 14px; }

        .entry-body { min-width: 0; }
        .entry-content { font-size: 14px; font-weight: 300; color: var(--text); line-height: 1.6; margin-bottom: 10px; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }
        .entry-meta { display: flex; align-items: center; gap: 10px; flex-wrap: wrap; }
        .entry-date { font-size: 11px; color: var(--muted); }
        .sentiment-badge { font-size: 10px; font-weight: 500; letter-spacing: 0.07em; text-transform: uppercase; padding: 2px 8px; border-radius: 20px; }
        .sentiment-badge.positive { background: rgba(109,191,138,0.12); color: #6dbf8a; }
        .sentiment-badge.negative { background: rgba(224,92,92,0.1); color: #e07070; }
        .sentiment-badge.neutral  { background: rgba(107,103,96,0.2); color: var(--muted2); }

        .entry-side { display: flex; align-items: flex-start; padding-top: 2px; }
        .btn-delete { width: 28px; height: 28px; border-radius: var(--r-sm); background: none; border: 1px solid transparent; color: var(--muted); font-size: 16px; cursor: pointer; display: flex; align-items: center; justify-content: center; transition: all 0.15s; line-height: 1; font-family: var(--sans); opacity: 0; }
        .entry-item:hover .btn-delete { opacity: 1; }
        .btn-delete:hover { background: rgba(224,92,92,0.1); border-color: rgba(224,92,92,0.2); color: var(--danger); }

        .empty-state { text-align: center; padding: 56px 20px; color: var(--muted); }
        .empty-state .big { font-family: var(--serif); font-size: 52px; color: var(--accent); opacity: 0.3; margin-bottom: 12px; }
        .empty-state p { font-size: 14px; }

        .loading { display: flex; align-items: center; justify-content: center; padding: 56px; color: var(--muted); font-size: 13px; gap: 10px; }
        .spinner { width: 16px; height: 16px; border: 1.5px solid var(--border2); border-top-color: var(--accent); border-radius: 50%; animation: spin 0.7s linear infinite; }
      `}</style>

      <aside className="sidebar">
        <div className="sidebar-logo">MoodLog</div>
        <div className="sidebar-link active">
          <span style={{fontSize:14}}>📓</span> Journal
        </div>
        <div className="sidebar-bottom">
          <button className="btn-logout" onClick={onLogout}>↩ Sign out</button>
        </div>
      </aside>

      <main className="main">
        <div className="page-header">
          <h1>{greeting}</h1>
          <p>{today}</p>
        </div>

        <div className="stats-row">
          <div className="stat-card">
            <div className="stat-label">Avg mood</div>
            <div className="stat-value" style={{ color: entries.length ? MOOD_COLOR[Math.round(parseFloat(avgMood(entries)))] : 'var(--accent)' }}>
              {avgMood(entries)}
            </div>
            <div className="stat-sub">out of 10</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">Entries</div>
            <div className="stat-value">{entries.length}</div>
            <div className="stat-sub">all time</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">Best day</div>
            <div className="stat-value" style={{ fontSize: 24, paddingTop: 4 }}>{bestDay(entries)}</div>
            <div className="stat-sub">by avg mood</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">Trend</div>
            <div className="stat-value" style={{ fontSize: trend ? 20 : 32, paddingTop: trend ? 6 : 0 }}>
              {trend ? (trend.dir === 'up' ? '↑' : trend.dir === 'down' ? '↓' : '→') : '—'}
            </div>
            {trend && <div className="stat-trend" style={{ color: trend.color }}>{trend.text}</div>}
            {!trend && <div className="stat-sub">need more data</div>}
          </div>
        </div>

        <InsightPanel entries={entries} />

        <div className="new-entry-card">
          <h2>New entry</h2>
          {formError && <div className="error-msg">{formError}</div>}
          <form onSubmit={submit}>
            <textarea
              ref={textareaRef}
              className="entry-textarea"
              placeholder="How are you feeling right now? What's on your mind?"
              value={content}
              onChange={e => setContent(e.target.value)}
              rows={4}
              maxLength={1000}
            />
            <div className="mood-row">
              <span className="mood-label">Mood</span>
              <input
                type="range" min={1} max={10} step={1}
                value={mood}
                onChange={e => setMood(Number(e.target.value))}
                className="mood-slider"
                style={{ background: `linear-gradient(to right, ${MOOD_COLOR[mood]} ${(mood-1)/9*100}%, var(--border) ${(mood-1)/9*100}%)` }}
              />
              <span className="mood-value" style={{ color: MOOD_COLOR[mood] }}>{mood}</span>
              <span style={{ fontSize: 20 }}>{MOOD_EMOJI[mood]}</span>
            </div>
            <div className="entry-actions">
              <span className="char-count">{content.length}/1000</span>
              <button className="btn-submit" type="submit" disabled={submitting || !content.trim()}>
                {submitting ? 'Saving...' : 'Save entry →'}
              </button>
            </div>
          </form>
        </div>

        <div className="entries-header">
          <h2>Past entries</h2>
          <span className="entries-count">{entries.length} total</span>
        </div>

        {loading ? (
          <div className="loading"><div className="spinner" /> Loading entries...</div>
        ) : entries.length === 0 ? (
          <div className="empty-state">
            <div className="big">✦</div>
            <p>No entries yet — write your first one above.</p>
          </div>
        ) : (
          entries.map(e => (
            <div className={`entry-item${deletingId === e.id ? ' deleting' : ''}`} key={e.id}>
              <div className="entry-mood">
                <span className="entry-mood-num" style={{ color: MOOD_COLOR[e.mood_score] }}>{e.mood_score}</span>
                <span className="entry-mood-emoji">{MOOD_EMOJI[e.mood_score]}</span>
              </div>
              <div className="entry-body">
                <div className="entry-content">{e.content}</div>
                <div className="entry-meta">
                  <span className="entry-date">{formatDate(e.created_at)}</span>
                  {e.sentiment && (
                    <span className={`sentiment-badge ${e.sentiment}`}>{e.sentiment}</span>
                  )}
                </div>
              </div>
              <div className="entry-side">
                <button
                  className="btn-delete"
                  onClick={() => remove(e.id)}
                  title="Delete entry"
                >
                  ✕
                </button>
              </div>
            </div>
          ))
        )}
      </main>
    </div>
  )
}