import React from 'react'
import { render } from '@testing-library/react';
import 'jest-styled-components';
import { Button } from '.';

test('Button works', () => {
  const { container } = render(<Button />)
  expect(container.firstChild).toMatchSnapshot();
  expect(container.firstChild).toHaveStyleRule('color','#fff');
})

test('Button works with props', () => {
  const { container } = render(<Button bgColor='black' />)
  expect(container.firstChild).toMatchSnapshot();
  expect(container.firstChild).toHaveStyleRule('background-color','black');
})
