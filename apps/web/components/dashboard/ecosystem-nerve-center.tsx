import React from 'react';
import { Activity, BarChart3, Database, Globe, Layers, ShieldCheck, Zap } from 'lucide-react';
import Link from 'next/link';

export function EcosystemNerveCenter() {
  return (
    <section className="py-12 bg-slate-950 text-white overflow-hidden">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col lg:flex-row gap-8">
          
          {/* Main Monitor: Trinity Balance */}
          <div className="flex-1 rounded-[40px] bg-gradient-to-br from-slate-900 to-slate-800 p-8 border border-white/10 relative group">
            <div className="absolute inset-0 bg-sky-500/5 opacity-0 group-hover:opacity-100 transition-opacity duration-700 blur-3xl" />
            <div className="relative z-10">
              <div className="flex items-center justify-between mb-8">
                <div>
                  <h3 className="text-xl font-black tracking-tight tracking-[-0.03em]">GLOBAL TRINITY MONITOR</h3>
                  <p className="text-slate-400 text-xs font-bold uppercase tracking-widest mt-1">Real-time Biosphere Balance</p>
                </div>
                <div className="h-10 w-10 rounded-full bg-sky-500/10 flex items-center justify-center ring-1 ring-sky-500/30">
                  <Activity className="h-5 w-5 text-sky-400 animate-pulse" />
                </div>
              </div>

              {/* Simulated Complex Radar/Chart Area */}
              <div className="grid grid-cols-3 gap-6 py-8">
                {[
                  { label: 'BODY', value: '92', color: 'bg-emerald-500', desc: 'Metabolic Sync' },
                  { label: 'ENV', value: '74', color: 'bg-sky-500', desc: 'Particulate Matter' },
                  { label: 'MIND', value: '88', color: 'bg-violet-500', desc: 'Alpha Wave Avg' },
                ].map((item) => (
                  <div key={item.label} className="text-center">
                    <div className="relative inline-flex items-center justify-center mb-4">
                      <svg className="h-24 w-24 -rotate-90">
                        <circle className="text-white/5" strokeWidth="6" stroke="currentColor" fill="transparent" r="44" cx="48" cy="48" />
                        <circle 
                          className={`transition-all duration-1000 ease-out`} 
                          strokeWidth="6" 
                          strokeDasharray={2 * Math.PI * 44} 
                          strokeDashoffset={2 * Math.PI * 44 * (1 - parseInt(item.value)/100)} 
                          strokeLinecap="round" 
                          stroke="currentColor" 
                          fill="transparent" 
                          r="44" cx="48" cy="48"
                          style={{ color: item.color.replace('bg-', 'var(--') }}
                        />
                      </svg>
                      <span className="absolute text-2xl font-black">{item.value}%</span>
                    </div>
                    <p className="text-sm font-bold text-white mb-0.5">{item.label}</p>
                    <p className="text-[10px] text-slate-500 font-bold uppercase tracking-tighter">{item.desc}</p>
                  </div>
                ))}
              </div>

              <div className="mt-8 pt-8 border-t border-white/5 flex items-center justify-between">
                <div className="flex -space-x-2">
                  {[1,2,3,4].map(i => (
                    <div key={i} className="h-8 w-8 rounded-full bg-slate-700 border-2 border-slate-900 ring-1 ring-white/10" />
                  ))}
                  <div className="h-8 w-8 rounded-full bg-sky-600 border-2 border-slate-900 flex items-center justify-center text-[10px] font-black">+124</div>
                </div>
                <div className="flex items-center gap-4">
                  <Link href="/domains/predictor" className="text-xs font-black uppercase tracking-widest text-sky-400 hover:text-sky-300 transition-colors flex items-center gap-2">
                    View Full Metrics <BarChart3 className="h-3 w-3" />
                  </Link>
                  <div className="h-4 w-px bg-white/10" />
                  <Link href="/domains/clinical" className="text-xs font-black uppercase tracking-widest text-emerald-400 hover:text-emerald-300 transition-colors flex items-center gap-2">
                    Launch Full Demo <Zap className="h-3 w-3" />
                  </Link>
                </div>
              </div>
            </div>
          </div>

          {/* Side Panel: System Pulse */}
          <div className="w-full lg:w-96 flex flex-col gap-6">
            
            {/* AI Agent Activity Log */}
            <div className="flex-1 rounded-[32px] bg-slate-900 p-6 border border-white/10 relative overflow-hidden">
               <h4 className="text-xs font-bold text-slate-500 uppercase tracking-widest mb-4 flex items-center gap-2">
                 <Zap className="h-3 w-3 text-amber-500" /> AI AGENT PIPELINE
               </h4>
               <div className="space-y-4">
                 {[
                   { id: 'clinical', domain: 'CLINICAL', msg: 'Packet hash-chain verified', time: 'now' },
                   { id: 'predictor', domain: 'PREDICTOR', msg: 'GNN-Relay weights updated', time: '2s' },
                   { id: 'reward', domain: 'REWARD', msg: 'Smart contract executed', time: '12s' },
                   { id: 'gxp', domain: 'GXP', msg: '電子署명(E-Sign) valid', time: '1m' },
                 ].map((log, i) => (
                   <Link href={`/domains/${log.id}`} key={i} className="flex gap-3 border-l border-white/5 pl-4 py-1 relative hover:bg-white/5 rounded-r-lg transition-colors group/log">
                     <div className="absolute -left-[1.5px] top-1/2 -translate-y-1/2 h-2 w-2 rounded-full bg-sky-500 group-hover/log:scale-125 transition-transform" />
                     <div className="flex-1">
                       <p className="text-[10px] text-slate-500 font-black tracking-tighter">{log.domain} ENGINE</p>
                       <p className="text-[13px] font-medium text-slate-300 leading-tight">{log.msg}</p>
                     </div>
                     <span className="text-[10px] text-slate-600 font-bold">{log.time}</span>
                   </Link>
                 ))}
               </div>
            </div>

            {/* Global Connectivity Status */}
            <Link href="/domains/hardware-core" className="block rounded-[32px] bg-sky-600 p-6 text-white relative group cursor-pointer overflow-hidden">
               <div className="absolute inset-0 bg-white/10 translate-x-full group-hover:translate-x-0 transition-transform duration-500" />
               <div className="relative z-10 flex items-center justify-between">
                 <div>
                   <p className="text-sky-200 text-[10px] font-black uppercase tracking-widest mb-1">Global Connect</p>
                   <h4 className="text-xl font-black">1.2M+ ACTIVE READERS</h4>
                 </div>
                 <Globe className="h-10 w-10 text-sky-200/50" />
               </div>
            </Link>

          </div>

        </div>
      </div>
    </section>
  );
}
