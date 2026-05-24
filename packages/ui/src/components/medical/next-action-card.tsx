import React from "react";

export interface NextAction {
  id: string;
  title: string;
  description: string;
  icon?: React.ReactNode;
  href?: string;
  onClick?: () => void;
  priority?: 'high' | 'medium' | 'low';
}

export interface NextActionCardProps {
  title: string;
  description: string;
  icon?: React.ReactNode;
  onClick?: () => void;
  priority?: 'high' | 'medium' | 'low';
}

export function NextActionCard({ title, description, icon, onClick, priority = 'medium' }: NextActionCardProps) {
  const borderColor = priority === 'high' ? 'border-sky-400 bg-sky-50/50' : priority === 'low' ? 'border-slate-200' : 'border-slate-200';
  return (
    <button
      onClick={onClick}
      className={`block w-full text-left p-4 rounded-xl border ${borderColor} hover:border-sky-300 hover:shadow-sm transition-all focus:outline-none focus:ring-2 focus:ring-sky-500/30 group`}
      aria-label={`다음 단계: ${title}`}
    >
      <div className="flex items-start gap-3">
        {icon && (
          <div className="h-9 w-9 rounded-lg bg-sky-50 flex items-center justify-center text-sky-600 group-hover:bg-sky-100 flex-shrink-0">
            {icon}
          </div>
        )}
        <div className="flex-1 min-w-0">
          <h3 className="font-bold text-sm text-slate-900 group-hover:text-sky-700 transition-colors">{title}</h3>
          <p className="text-slate-500 text-xs mt-0.5 leading-relaxed">{description}</p>
        </div>
        <span className="text-slate-300 group-hover:text-sky-500 transition-colors flex-shrink-0 mt-0.5">→</span>
      </div>
    </button>
  );
}

export function NextActionGroup({ title, actions }: { title?: string; actions: NextAction[] }) {
  return (
    <section aria-label={title || '다음 단계'} className="space-y-2">
      {title && (
        <h3 className="text-sm font-bold text-slate-500 mb-3 flex items-center gap-1.5">
          <span className="h-1 w-1 rounded-full bg-sky-500" />
          {title}
        </h3>
      )}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
        {actions.map(action => (
          <div key={action.id}>
            <NextActionCard
              title={action.title}
              description={action.description}
              icon={action.icon}
              onClick={action.onClick}
              priority={action.priority}
            />
          </div>
        ))}
      </div>
    </section>
  );
}
