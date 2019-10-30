import React from 'react';
import styled from 'styled-components';

const TextFieldContainer = styled.div`
  display: flex;
  flex-direction: column;
  font-family: 'Manulife JH Sans', sans-serif;
  position: relative;
  width: 100%;
  @media only screen and (max-width : 414px) {
    margin-bottom: 1rem;
  }
`;

const Label = styled.label`
  color: #8e90a2;
  cursor: text;
  display: flex;
  flex-direction: column;
  font-size: 0.875rem;
  padding: 0;
  position: relative;
`;

const Input = styled.input`
  background: none;
  border: none;
  box-shadow: 0 1px 0 0 #8e90a2;
  color: #282b3e;
  font-size: 1rem;
  line-height: 1.875rem;
  outline: none;
  padding: 0.3125rem 0;
  transition: box-shadow 240ms ease-in-out;
  font-family: 'Manulife JH Sans', sans-serif;

  &:focus {
    box-shadow: 0 2px 0 0 #00a758;
  }
`;

const Error = styled.div`
  color: #dc5a44;
  display: flex;
  font-size: 0.875rem;
  position: absolute;
  bottom: calc(-10px + -0.875rem);
  @media only screen and (max-width : 414px) {
    bottom: calc(-30px + -0.875rem);
  }
  `;

export const TextField = ({ inputLabel, name, onChange, placeholder, readOnly, useRef, value, error }) => (
  <TextFieldContainer>
    <Label htmlFor={name}>{inputLabel}</Label>
    <Input
      type="text"
      name={name}
      id={name}
      placeholder={placeholder}
      onChange={onChange}
      value={value}
      readOnly={readOnly}
      ref={useRef}
      />
    {error && <Error role="alert">{error}</Error>}
  </TextFieldContainer>
)
