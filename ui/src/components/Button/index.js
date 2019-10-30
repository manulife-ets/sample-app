import styled from 'styled-components';

export const Button = styled.button`
  min-width: 150px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0;
  text-align: center;
  vertical-align: middle;
  touch-action: manipulation;
  cursor: pointer;
  border: 1px solid transparent;
  white-space: nowrap;
  user-select: none;
  font-weight: 600;
  max-width: 100%;
  white-space: normal;
  background-color: ${props => props.bgColor || "#ec6453"};
  color: #fff;
  
  @media only screen and (max-width : 414px) {
    padding: 1rem 0;
  }

  &:hover {
    background-color: ${props => props.hoverColor || "#dc5a44"}
  }
`;