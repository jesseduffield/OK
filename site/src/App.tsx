import React from 'react';
import './App.css';
import { templates } from './templates';
import { Controlled as CodeMirror } from 'react-codemirror2';
import { contentToUrl, urlToContent } from './utils/encoding';
import config from './config';

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

const App = () => {
  const [value, setValue] = React.useState(getInitialValue);
  const [output, setOutput] = React.useState('');
  const [runComplete, setRunComplete] = React.useState(false);

  const onRun = async () => {
    setOutput('Running...');
    setRunComplete(false);
    const output = await run(value);
    setOutput(output);
    setRunComplete(true);
  };

  return (
    <div className="App">
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
      </header>
      <Editor value={value} onNewValue={value => setValue(value)} />
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
}: {
  value: string;
  onNewValue: (value: string) => void;
}) => {
  return (
    <div className="Editor">
      <CodeMirror
        value={value}
        options={{ mode: 'go', theme: 'none', lineNumbers: true }}
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
