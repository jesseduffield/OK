import React from 'react';
import './App.css';
import { templates } from './templates';
import { Controlled as CodeMirror } from 'react-codemirror2';
import { contentToUrl, urlToContent } from './utils/encoding';
import config from './config';

import './ok';

const run = async (text: string) => {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    mode: 'cors',
    redirect: 'follow',
    body: JSON.stringify({ content: text }),
  };

  try {
    const response = await fetch(
      `${config.serverUrl}/run`,
      requestOptions as any
    );
    return await response.text();
  } catch (error) {
    alert(error);
    return 'An error occured';
  }
};

const getInitialValue = () =>
  urlToContent(window.location.href) || templates[0].value;

const eighteenYearsAgo = new Date();
eighteenYearsAgo.setFullYear(eighteenYearsAgo.getFullYear() - 18);

const HighlightingModal = ({
  onClose,
  onConfirm,
}: {
  onClose: () => void;
  onConfirm: () => void;
}) => {
  const [date, setDate] = React.useState<null | Date>(null);

  const canEnable = date && date > eighteenYearsAgo;

  return (
    <div className="Modal-outer" onClick={onClose}>
      <div className="Modal-middle">
        <div className="Modal-inner" onClick={event => event.stopPropagation()}>
          <div>
            <h1 className="Modal-heading">WARNING: GO BACK!</h1>
            <p>Syntax highlighting is juvenile and you don't need it.</p>
            <p>
              Please enter your date of birth to confirm that you are{' '}
              <span className="italic">under</span> 18:
            </p>
            <input
              type="date"
              onChange={event => {
                setDate(new Date(event.target.value));
              }}
            ></input>
          </div>
          <div className="Modal-buttons">
            <button
              className={`enable-button ${!canEnable && 'disabled'}`}
              onClick={onConfirm}
              disabled={!canEnable}
            >
              Enable
            </button>
            <Button onClick={onClose}>Cancel</Button>
          </div>
        </div>
      </div>
    </div>
  );
};

const App = () => {
  const [value, setValue] = React.useState(getInitialValue);
  const [output, setOutput] = React.useState('');
  const [runComplete, setRunComplete] = React.useState(false);
  const [codeMirrorTheme, setCodeMirrorTheme] = React.useState('');
  const [highlightingModalOpen, setHighlightingModalOpen] = React.useState(
    false
  );

  const enableHighlighting = () => {
    setCodeMirrorTheme('default');
    setHighlightingModalOpen(false);
  };

  const onRun = async () => {
    setOutput('Running...');
    setRunComplete(false);
    const output = await run(value);
    setOutput(output);
    setRunComplete(true);
  };

  return (
    <div className="App">
      {highlightingModalOpen && (
        <HighlightingModal
          onClose={() => setHighlightingModalOpen(false)}
          onConfirm={enableHighlighting}
        />
      )}
      <header className="Banner">
        <div className="Head">
          The <span className="OK">OK?</span> Playground
        </div>
        <Button onClick={onRun}>Run</Button>
        <Share value={value} />
        <select
          className="Select"
          onChange={event => setValue(event.target.value)}
        >
          {templates.map(({ label, value }) => (
            <option key={label} value={value}>
              {label}
            </option>
          ))}
        </select>
        <Button
          onClick={() => {
            window!.open(config.docsUrl, '_blank')!.focus();
          }}
        >
          Docs
        </Button>
        <Button
          onClick={() => {
            setHighlightingModalOpen(true);
          }}
        >
          Enable Syntax Highlighting
        </Button>
      </header>
      <Editor
        value={value}
        onNewValue={value => setValue(value)}
        theme={codeMirrorTheme}
      />
      <p className="Output">{output}</p>
      {runComplete && (
        <p className="Output RunComplete">Program has finished.</p>
      )}
    </div>
  );
};

const Button = ({
  children,
  onClick,
}: {
  children: string;
  onClick: () => void;
}) => (
  <button className="Button" onClick={onClick}>
    {children}
  </button>
);

const Editor = ({
  value,
  onNewValue,
  theme,
}: {
  value: string;
  onNewValue: (value: string) => void;
  theme: string;
}) => {
  return (
    <div className="Editor">
      <CodeMirror
        value={value}
        options={{ mode: 'ok', theme, lineNumbers: true }}
        onBeforeChange={(_editor, _data, value) => {
          onNewValue(value);
        }}
      />
    </div>
  );
};

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text);
};

const CopyButton = ({ text }: { text: string }) => {
  const [copied, setCopied] = React.useState(false);

  return (
    <Button
      onClick={() => {
        copyToClipboard(text);
        setCopied(true);
        setTimeout(() => setCopied(false), 1000);
      }}
    >
      {copied ? 'Copied' : 'Copy Text'}
    </Button>
  );
};

const Share = ({ value }: { value: string }) => {
  const [shareUrl, setShareUrl] = React.useState<string | null>(null);

  return (
    <>
      <Button onClick={() => setShareUrl(contentToUrl(config.baseUrl, value))}>
        Share
      </Button>
      {shareUrl !== null && (
        <>
          <div className="ShareInput">{shareUrl}</div>
          <CopyButton text={shareUrl} />
        </>
      )}
    </>
  );
};

export default App;
