import React, {useEffect, useState} from 'react'

export default function App(){
  const [status,setStatus]=useState({});
  const [loading,setLoading]=useState(false);
  const [backend,setBackend]=useState('');
  const [query,setQuery]=useState('SELECT 1');
  const [result,setResult]=useState(null);
  const [error,setError]=useState(null);

  const fetchPing=async ()=>{
    setLoading(true); setError(null);
    try{
      const r=await fetch('/ping');
      const j=await r.json(); setStatus(j);
      const keys=Object.keys(j); if(keys.length && !backend) setBackend(keys[0]);
    }catch(e){ setError('failed to fetch /ping: '+e.message) }
    setLoading(false);
  }

  useEffect(()=>{ fetchPing() },[])

  const doQuery=async (ev)=>{
    ev && ev.preventDefault(); setError(null); setResult(null); setLoading(true);
    if(!backend){ setError('choose backend'); setLoading(false); return }
    try{
      const r=await fetch('/query',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({backend,query})});
      const j=await r.json();
      if(!r.ok){ setError(j.error||JSON.stringify(j)) } else { setResult(j.rows) }
    }catch(e){ setError('request failed: '+e.message) }
    setLoading(false);
  }

  return (
    <div className="container my-4">
      <div className="d-flex justify-content-between align-items-center mb-3">
        <h3>db-connect-demo</h3>
        <div>
          <button className="btn btn-sm btn-outline-primary me-2" onClick={fetchPing} disabled={loading}>Refresh</button>
        </div>
      </div>
      <div className="row">
        <div className="col-md-5">
          <div className="card">
            <div className="card-header">Backends / Health</div>
            <div className="card-body" style={{maxHeight:400, overflow:'auto'}}>
              {Object.keys(status).length===0 && <div className="text-muted">No backends registered.</div>}
              {Object.entries(status).map(([k,v])=> (
                <div key={k} className="mb-2 p-2 border rounded">
                  <div className="d-flex justify-content-between"><strong>{k}</strong><small>{v==='ok'?<span className='text-success'>ok</span>:<span className='text-danger'>{v}</span>}</small></div>
                </div>
              ))}
            </div>
          </div>
        </div>
        <div className="col-md-7">
          <div className="card">
            <div className="card-header">Query</div>
            <div className="card-body">
              <form onSubmit={doQuery}>
                <div className="mb-2">
                  <label className="form-label">Backend</label>
                  <select className="form-select" value={backend} onChange={e=>setBackend(e.target.value)}>
                    {Object.keys(status).map(k=> <option key={k} value={k}>{k}</option>)}
                  </select>
                </div>
                <div className="mb-2">
                  <label className="form-label">Query</label>
                  <textarea className="form-control" rows={4} value={query} onChange={e=>setQuery(e.target.value)} />
                </div>
                <div className="mb-2">
                  <button className="btn btn-primary me-2" type="submit" disabled={loading}>Execute</button>
                  <button type="button" className="btn btn-outline-secondary" onClick={()=>setQuery('SELECT 1')}>Reset</button>
                </div>
              </form>
              {error && <div className="alert alert-danger mt-2">{error}</div>}
              {result && <div className="mt-2"><h6>Result</h6><pre className="p-2 bg-light border">{JSON.stringify(result,null,2)}</pre></div>}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
