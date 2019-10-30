import React, { useState } from 'react';
import Grid from '@radd/react/dist/components/containers/Grid';
import { TextField } from '../TextField';
import { Button } from '../Button';

const fetcher = async ({url, method, body}) => {
  return fetch(url, {
    method: method,
    headers: {
      'Content-Type': 'application/json',
    },
    body: body
  });
}

const ShortenerForm = ({task}) => {
  const [longUrl, setLongUrl] = useState('');
  const [urlId, setUrlId] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [errorMsg, setErrorMsg] = useState('');

  const updateUrl = e => {
    setLongUrl(e.target.value);
    setErrorMsg(false);
  }

  const updateUrlId = e => {
    setUrlId(e.target.value);
  }

  const handleSubmit = e => {
    e.preventDefault();
    processShortenRequest(longUrl);
  }

  const processShortenRequest = async str => {
    const url = isValidUrl(str) && str;
    if (str.length < 1) {
      renderError('EMPTY_INPUT');
      return false;
    }
    if (!url) {
      renderError('INVALID_INPUT');
      return false;
    } else {
      await makeShortenUrlRequest(url);
    }
  }

  const isValidUrl = str => {
    const urlPattern = new RegExp(/^https?:\/\/(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[-A-Z0-9+&@#\/%=~_|$?!:,.])*(?:\([-A-Z0-9+&@#\/%=~_|$?!:,.]*\)|[A-Z0-9+&@#\/%=~_|$])$/, 'i');
    return !!urlPattern.test(str);
  }

  const buildRequestFromTask = (task, url) => {
    switch (task) {
      case 'SHORTEN':
        return {
          url: '/api/v2',
          method: 'POST',
          body: JSON.stringify({
            longURL: url
          })
        }
      case 'UPDATE':
        return {
          url: '/api/v2/' + urlId,
          method: 'PUT',
          body: JSON.stringify({
            longURL: url
          })
        }
      case 'CREATE':
        return {
          url: '/api/v2',
          method: 'POST',
          body: JSON.stringify({
            id: urlId,
            longURL: url
          })
        }
      default: renderError('EMPTY')
    }

  }

  const makeShortenUrlRequest = async url => {
    let shortUrl;
    const requestParams = buildRequestFromTask(task, url);
    try {
      let response = await fetcher(requestParams);
      let json = await response.json()
      if (!response.ok) {
        renderError('API_ERROR', json.message)
        return false;
      } else {
        shortUrl = json.id;
        setShortUrl(shortUrl);
        document.getElementById('shortUrl').select();
      }
    } catch(err) {
      renderError('API_ERROR', 'URL ID not found');
    }
  }

  const handleCopy = e => {
    e.preventDefault();
    copyShortUrl();
  }

  const copyShortUrl = () => {
    const shortUrl = document.getElementById('shortUrl');
    shortUrl.select();
    document.execCommand('copy');
  }

  const renderError = (type, message = 'Error calling API') => {
    const errorTypes = new Map([
      ['EMPTY_INPUT', 'Please enter a URL'],
      ['INVALID_INPUT', 'URL is not valid'],
      ['API_ERROR', message],
      ['EMPTY', '&nbsp;']
    ])
    setErrorMsg(errorTypes.get(type));
  }

  const placeholder = task === 'CREATE' ? 'Like a-custom-link-name or ✔🙌👍' : 'ex: custom-link-name or oqw3k78oz51m83y';
  const slogans = ['do-the-right-thing', 'think-big', 'get-it-done-together', 'own-it', 'share-your-humanity', 'obsess-about-customers'];
  const slogan = slogans[Math.floor(Math.random() * slogans.length)];
  // use dev api if running locally
  const apiUrl = window.location.origin.indexOf('localhost') > 0 ? 'https://go-dev.manulife.com' : window.location.origin;
  return (
  <>
    <form onSubmit={handleSubmit} autoComplete="off">
        {((task === 'CREATE') || (task === 'UPDATE')) &&
        <Grid xs="3fr 1fr" gutter="0" margin="2rem none 2rem">
          <TextField
            value={urlId}
            inputLabel="Enter a URL ID"
            name="urlId"
            placeholder={placeholder}
            onChange={updateUrlId}
          />
        </Grid>
      }
      <Grid xs="3fr 1fr" gutter="0" margin="none none 2rem">
        <TextField
          value={longUrl}
          inputLabel="Enter a looooong URL"
          name="longUrl"
            onChange={updateUrl}
            placeholder={`ex: https://www.manulife.com/really-really-really-really-really-long-url-to-${slogan}`}
          error={errorMsg}
        />
        <Button>{task.charAt(0).toUpperCase() + task.slice(1).toLowerCase()} It!</Button>
      </Grid>
    </form>
    {shortUrl &&
      <form onSubmit={handleCopy}>
        <Grid xs="3fr 1fr" gutter="0">
          <TextField
            value={`${apiUrl}/${shortUrl}`}
            inputLabel="Now copy that URL"
            name="shortUrl"
            readOnly={true}
          />
          <Button bgColor="#06c7ba" hoverColor="#08a298">Copy It!</Button>
        </Grid>
      </form>
    }
  </>
  );
}

export default ShortenerForm;
