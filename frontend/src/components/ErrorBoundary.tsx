import { Component, type ErrorInfo, type ReactNode } from 'react';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
}

// Catches render-time throws so a bad API payload can't blank the whole app.
export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false };

  static getDerivedStateFromError(): State {
    return { hasError: true };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('Uncaught render error:', error, info);
  }

  render() {
    if (!this.state.hasError) {
      return this.props.children;
    }

    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <div className="card max-w-md text-center">
          <div className="mb-4 text-3xl text-error">!</div>
          <h2 className="mb-2 font-serif text-lg font-bold text-ink">Something went wrong</h2>
          <p className="mb-4 text-sm text-feather">
            The page hit an unexpected error. Reloading usually sorts it.
          </p>
          <button className="btn-primary" onClick={() => window.location.reload()}>
            Reload
          </button>
        </div>
      </div>
    );
  }
}
