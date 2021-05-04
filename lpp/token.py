from enum import (
    auto,
    Enum,
    unique
)

from typing import NamedTuple, Dict

@unique
class TokenType(Enum):
    ASSING = auto()
    COMMA = auto()
    EOF = auto()
    FUNCTION = auto()
    IDENT = auto()
    ILLEGAL = auto()
    INT = auto()
    LBRACE = auto()
    LET = auto()
    LPAREN = auto()
    PLUS = auto()
    RBRACE = auto()
    RPAREN = auto()
    SEMICOLON = auto()

class Token(NamedTuple):
    """Token definition"""

    token_type: TokenType
    literal: str

    def __str__(self) -> str:
        return f'Type: {self.token_type}, Literal: {self.literal}'

def lookup_token_type(literal: str) -> TokenType:
    """check if the given literal is a keyword"""
    
    keywords: Dict[str, TokenType] = {
        'var': TokenType.LET,
        'funcion': TokenType.FUNCTION
    }

    return keywords.get(literal, TokenType.IDENT)
