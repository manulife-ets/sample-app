import React from 'react'
import { render } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect';
import 'jest-styled-components';
import { TextField } from '.';

describe('Textfield', () => {
  test('it works', () => {
    const { container } = render(<TextField />);
    expect(container.firstChild).toMatchSnapshot();
    expect(container.firstChild).toHaveStyleRule('font-family', "'Manulife JH Sans',sans-serif");
    
  })

  test('it renders an error', async () => {
    const { findByRole } = render(<TextField error='BOOM' />);
    const alert = await findByRole('alert');
    expect(alert).toHaveTextContent(/BOOM/i)
  })
})