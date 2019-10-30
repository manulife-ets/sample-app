import React from 'react'
import { render, fireEvent, findByText } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import ShortenerForm from '.';
import { act } from 'react-dom/test-utils';

describe('ShortenerForm', () => {
  test('it works', () => {
    const { container } = render(<ShortenerForm task='TEST' />);
    expect(container.firstChild).toMatchSnapshot();
  });

  test('it validates the URL', async () => {
    const { container, getByLabelText, getByText, findByRole } = render(<ShortenerForm task='TEST' />);
    const handleSubmit = jest.fn()

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: 'htw/:qw.c/' } });
    fireEvent.click(getByText(/test it/i));

    const alert = await findByRole('alert');
    expect(alert).toHaveTextContent(/url is not valid/i);
  });

  test('it requires a URL to be entered', async () => {
    const { container, getByLabelText, getByText, findByRole } = render(<ShortenerForm task='TEST' />);

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: '' } });
    fireEvent.click(getByText(/test it/i));

    const alert = await findByRole('alert');
    expect(alert).toHaveTextContent(/please enter a url/i);
  });

  test('it clears the error message when the input changes', async () => {
    const { container, getByLabelText, getByText, findByRole } = render(<ShortenerForm task='TEST' />);

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: 'errrorrr cawzer' } });
    fireEvent.click(getByText(/test it/i));

    let alert = await findByRole('alert');
    expect(alert).toBeInTheDocument();

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: 'errrorrr caw' } });

    alert = await document.querySelector('[role="alert"]');
    expect(alert).not.toBeInTheDocument();

  });

  test('it allows a valid URL', async () => {
    // mock out window.fetch for the test
    const makeShortenUrlRequest = jest.fn();
    const { container, getByLabelText, getByText } = render(<ShortenerForm task='SHORTEN' />);

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: 'http://test.com' } });
    fireEvent.click(getByText(/shorten it/i));

    alert = await document.querySelector('[role="alert"]');
    expect(alert).not.toBeInTheDocument();
  });

  test('it renders the correct CREATE form', async () => {
    // mock out window.fetch for the test
    const makeShortenUrlRequest = jest.fn();
    const { container, getByLabelText, getByText } = render(<ShortenerForm task='CREATE' />);

    expect(getByText(/create it/i)).toBeInTheDocument();
  });

  test('it renders the correct CREATE form', async () => {
    // mock out window.fetch for the test
    const makeShortenUrlRequest = jest.fn();
    const { container, getByLabelText, getByText } = render(<ShortenerForm task='CREATE' />);

    expect(getByText(/create it/i)).toBeInTheDocument();
  });

  test('it renders the correct UPDATE form', async () => {
    // mock out window.fetch for the test
    const makeShortenUrlRequest = jest.fn();
    const { container, getByLabelText, getByText } = render(<ShortenerForm task='UPDATE' />);

    expect(getByText(/update it/i)).toBeInTheDocument();
  });

  test('it renders the API error message', async () => {
    global.fetch = jest.fn().mockImplementationOnce(() => {
      return new Promise((resolve, reject) => {
        resolve({
          ok: false,
          status: 400,
          json: () => { return { id: '🃏' } }
        });
      });
    });

    const { container, getByLabelText, getByText, findByRole } = render(<ShortenerForm task='TEST' />);

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: 'errrorrr cawzer' } });
    fireEvent.click(getByText(/test it/i));

    const alert = await findByRole('alert');
    expect(alert).toBeInTheDocument();

  });

  test('it renders the copy URL form from the SHORTEN API response', async () => {
    global.fetch = jest.fn().mockImplementationOnce(() => {
      return new Promise((resolve, reject) => {
        resolve({
          ok: true,
          status: 200,
          json: () => { return { id: 'JESTTEST' } }
        });
      });
    });

    const { container, getByLabelText, getByText, findByText } = render(<ShortenerForm task='SHORTEN' />);

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: 'http://www.com' } });
    fireEvent.click(getByText(/shorten it/i));

    const copy = await findByText('Copy It!');
    expect(copy).toBeInTheDocument();

    const input = getByLabelText('Now copy that URL');
    expect(input.value).toEqual('https://go.manulife.com/JESTTEST');
  });

  test('it renders the copy URL form from the CREATE API response', async () => {
    global.fetch = jest.fn().mockImplementationOnce(() => {
      return new Promise((resolve, reject) => {
        resolve({
          ok: true,
          status: 200,
          json: () => { return { id: 'JESTTEST' } }
        });
      });
    });

    const { container, getByLabelText, getByText, findByText } = render(<ShortenerForm task='CREATE' />);

    fireEvent.change(getByLabelText(/enter a url id/i), { target: { value: '🃏' } });

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: 'http://www.com' } });
    fireEvent.click(getByText(/create it/i));

    const copy = await findByText('Copy It!');
    expect(copy).toBeInTheDocument();

    const input = getByLabelText('Now copy that URL');
    expect(input.value).toEqual('https://go.manulife.com/JESTTEST');
  });

  test('it renders the copy URL form from the CREATE API response', async () => {
    global.fetch = jest.fn().mockImplementationOnce(() => {
      return new Promise((resolve, reject) => {
        resolve({
          ok: true,
          status: 200,
          json: () => { return { id: '🃏' } }
        });
      });
    });

    const { container, getByLabelText, getByText, findByText } = render(<ShortenerForm task='UPDATE' />);

    fireEvent.change(getByLabelText(/enter a url id/i), { target: { value: '🃏' } });

    fireEvent.change(getByLabelText(/enter a lo+ng url/i), { target: { value: 'http://www.com' } });
    fireEvent.click(getByText(/update it/i));

    const copy = await findByText('Copy It!');
    expect(copy).toBeInTheDocument();

    const input = getByLabelText('Now copy that URL');
    expect(input.value).toEqual('https://go.manulife.com/🃏');
  });
});
