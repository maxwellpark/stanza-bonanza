import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { PoemFormat, ApprovalMode } from '@/types/poem';
import { useCreatePoem } from '@/hooks/usePoems';
import { formatPoemFormat } from '@/lib/utils';

var allFormats: PoemFormat[] = [
  'free_verse',
  'haiku',
  'sonnet',
  'limerick',
  'iambic_pentameter',
  'rhyming_couplets',
  'custom',
];

var approvalOptions: { value: ApprovalMode; label: string; description: string }[] = [
  { value: 'open', label: 'Open', description: 'Anyone can add stanzas immediately' },
  { value: 'approval_required', label: 'Requires Approval', description: 'Stanzas must be approved by you' },
  { value: 'closed', label: 'Closed', description: 'No one else can add stanzas' },
];

export function PoemEditor() {
  var navigate = useNavigate();
  var mutation = useCreatePoem();

  var [title, setTitle] = useState('');
  var [description, setDescription] = useState('');
  var [text, setText] = useState('');
  var [format, setFormat] = useState<PoemFormat>('free_verse');
  var [approvalMode, setApprovalMode] = useState<ApprovalMode>('open');
  var [maxStanzas, setMaxStanzas] = useState('');
  var [tags, setTags] = useState('');
  var [error, setError] = useState('');

  async function submit(draft: boolean) {
    setError('');

    if (!title.trim() || !text.trim()) {
      setError('Title and first stanza are required.');
      return;
    }

    var parsedTags = tags
      .split(',')
      .map((t) => t.trim())
      .filter(Boolean)
      .slice(0, 10);

    try {
      var poem = await mutation.mutateAsync({
        title: title.trim(),
        description: description.trim() || undefined,
        text: text.trim(),
        format,
        approvalMode,
        maxStanzas: maxStanzas ? Number(maxStanzas) : undefined,
        tags: parsedTags.length ? parsedTags : undefined,
        draft,
      });
      navigate(`/poems/${poem.id}`);
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Something went wrong. Please try again.');
      }
    }
  }

  return (
    <div className="mx-auto max-w-2xl">
      <div className="card">
        <div className="mb-6 border-b border-parchment-dark pb-4">
          <h1 className="font-serif text-2xl font-bold text-ink">Create a New Poem</h1>
          <p className="mt-1 font-sans text-sm text-feather">
            Start a poem and invite others to continue it.
          </p>
        </div>

        <form onSubmit={(e) => { e.preventDefault(); submit(false); }} className="flex flex-col gap-5">
          <div>
            <label htmlFor="title" className="mb-1 block font-sans text-sm font-medium text-ink">
              Title
            </label>
            <input
              id="title"
              type="text"
              required
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Give your poem a title..."
              className="w-full rounded-lg border border-parchment-dark bg-white px-4 py-2 font-serif text-lg text-ink outline-none transition-colors focus:border-accent"
            />
          </div>

          <div>
            <label htmlFor="description" className="mb-1 block font-sans text-sm font-medium text-ink">
              Description <span className="font-normal text-feather">(optional)</span>
            </label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Set the scene or give contributors guidance..."
              rows={2}
              className="w-full resize-none rounded-lg border border-parchment-dark bg-white px-4 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
            />
          </div>

          <div>
            <label htmlFor="text" className="mb-1 block font-sans text-sm font-medium text-ink">
              First Stanza
            </label>
            <textarea
              id="text"
              required
              value={text}
              onChange={(e) => setText(e.target.value)}
              placeholder="Write the opening stanza..."
              rows={6}
              className="w-full resize-none rounded-lg border border-parchment-dark bg-white px-4 py-3 text-base leading-relaxed text-ink outline-none transition-colors focus:border-accent"
              style={{ fontFamily: 'var(--font-body)' }}
            />
            {text.trim() && (
              <div className="mt-2">
                <span className="mb-1 block font-sans text-xs text-feather">Preview</span>
                <div className="prose prose-sm max-w-none rounded-lg border border-dashed border-parchment-dark bg-parchment/50 px-4 py-3 text-base leading-relaxed text-ink" style={{ fontFamily: 'var(--font-body)' }}>
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>{text.replace(/\n/g, '  \n')}</ReactMarkdown>
                </div>
              </div>
            )}
          </div>

          <div>
            <label htmlFor="format" className="mb-1 block font-sans text-sm font-medium text-ink">
              Format
            </label>
            <select
              id="format"
              value={format}
              onChange={(e) => setFormat(e.target.value as PoemFormat)}
              className="w-full rounded-lg border border-parchment-dark bg-white px-4 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
            >
              {allFormats.map((f) => (
                <option key={f} value={f}>{formatPoemFormat(f)}</option>
              ))}
            </select>
          </div>

          <div>
            <span className="mb-2 block font-sans text-sm font-medium text-ink">
              Collaboration Mode
            </span>
            <div className="flex flex-col gap-2">
              {approvalOptions.map((opt) => (
                <label key={opt.value} className="flex cursor-pointer items-start gap-3 rounded-lg border border-parchment-dark px-4 py-3 transition-colors hover:border-accent">
                  <input
                    type="radio"
                    name="approvalMode"
                    value={opt.value}
                    checked={approvalMode === opt.value}
                    onChange={() => setApprovalMode(opt.value)}
                    className="mt-0.5"
                  />
                  <div>
                    <span className="block font-sans text-sm font-medium text-ink">{opt.label}</span>
                    <span className="block font-sans text-xs text-feather">{opt.description}</span>
                  </div>
                </label>
              ))}
            </div>
          </div>

          <div>
            <label htmlFor="tags" className="mb-1 block font-sans text-sm font-medium text-ink">
              Tags <span className="font-normal text-feather">(optional, comma-separated)</span>
            </label>
            <input
              id="tags"
              type="text"
              value={tags}
              onChange={(e) => setTags(e.target.value)}
              placeholder="nature, love, autumn"
              className="w-full rounded-lg border border-parchment-dark bg-white px-4 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
            />
          </div>

          <div>
            <label htmlFor="maxStanzas" className="mb-1 block font-sans text-sm font-medium text-ink">
              Max Stanzas <span className="font-normal text-feather">(optional)</span>
            </label>
            <input
              id="maxStanzas"
              type="number"
              min={1}
              value={maxStanzas}
              onChange={(e) => setMaxStanzas(e.target.value)}
              placeholder="No limit"
              className="w-full rounded-lg border border-parchment-dark bg-white px-4 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
            />
          </div>

          {error && (
            <p className="font-sans text-sm text-error">{error}</p>
          )}

          <div className="flex gap-3">
            <button
              type="submit"
              disabled={mutation.isPending}
              className="btn-primary disabled:opacity-50"
            >
              {mutation.isPending ? 'Saving...' : 'Publish Poem'}
            </button>
            <button
              type="button"
              onClick={() => submit(true)}
              disabled={mutation.isPending}
              className="btn-secondary disabled:opacity-50"
            >
              Save as Draft
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
