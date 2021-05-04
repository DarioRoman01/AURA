from unittest import TestCase
from typing import List
from lpp.token import Token, TokenType
from lpp.lexer import Lexer

class LexerTest(TestCase):
    """
    Lexer tests cases verify that the lexer 
    assing tokens type propertly
    """

    def test_illegal(self) -> None:
        """Test illegal tokens."""

        source: str = '¡¿@'
        lexer: Lexer = Lexer(source)

        tokens: List[Token] = []
        for i in range(len(source)):
            tokens.append(lexer.next_token())

        expected_tokens: List[Token] = [
            Token(TokenType.ILLEGAL, '¡'),
            Token(TokenType.ILLEGAL, '¿'),
            Token(TokenType.ILLEGAL, '@'),
        ]

        self.assertEquals(tokens, expected_tokens)

    def test_one_character_operator(self) -> None:
        """Test operators character."""

        source: str = '=+'
        lexer: Lexer = Lexer(source)

        tokens: List[Token] = []
        for i in range(len(source)):
            tokens.append(lexer.next_token())
        
        expected_tokens: List[Token] = [
            Token(TokenType.ASSING, '='),
            Token(TokenType.PLUS, '+'),
        ]

        self.assertEquals(tokens, expected_tokens)

    def test_eof(self) -> None:
        """Test end of file tokens."""

        source: str = '+'
        lexer: Lexer = Lexer(source)
        tokens: List[Token] = []

        for i in range(len(source) + 1):
                tokens.append(lexer.next_token())
        
        expected_tokens: List[Token] = [
            Token(TokenType.PLUS, '+'),
            Token(TokenType.EOF, '')
        ]

        self.assertEqual(tokens, expected_tokens)

    def test_delimiters(self) -> None:
        """test delimiters tokens."""

        source = '(){},;'
        lexer: Lexer = Lexer(source)
        tokens: List[Token] = []

        for i in range(len(source)):
            tokens.append(lexer.next_token())

        expected_tokens: List[Token] = [
            Token(TokenType.LPAREN, '('),
            Token(TokenType.RPAREN, ')'),
            Token(TokenType.LBRACE, '{'),
            Token(TokenType.RBRACE, '}'),
            Token(TokenType.COMMA, ','),
            Token(TokenType.SEMICOLON, ';'),
        ]

        self.assertEquals(tokens, expected_tokens)
        
    def test_assingment(self) -> None:
        """test variables assingment tokens"""

        source: str = 'var cinco = 5;'
        lexer: Lexer = Lexer(source)

        tokens: List[Token] = []
        for i in range(5):
            tokens.append(lexer.next_token())
        
        expected_tokens: List[Token] = [
            Token(TokenType.LET, 'var'),
            Token(TokenType.IDENT, 'cinco'),
            Token(TokenType.ASSING, '='),
            Token(TokenType.INT, '5'),
            Token(TokenType.SEMICOLON, ';'),
        ]

        self.assertEquals(tokens, expected_tokens)

    def test_function_declaration(self) -> None:
        """Test functions declarations tokens."""

        source: str = '''
            var suma = funcion(x, y) {
                x + y;
            };
        '''

        lexer: Lexer = Lexer(source)
        tokens: List[Token] = []
        for i in range(16):
            tokens.append(lexer.next_token())

        expected_tokens: List[Token] = [
            Token(TokenType.LET, 'var'),
            Token(TokenType.IDENT, 'suma'),
            Token(TokenType.ASSING, '='),
            Token(TokenType.FUNCTION, 'funcion'),
            Token(TokenType.LPAREN, '('),
            Token(TokenType.IDENT, 'x'),
            Token(TokenType.COMMA, ','),
            Token(TokenType.IDENT, 'y'),
            Token(TokenType.RPAREN, ')'),
            Token(TokenType.LBRACE, '{'),
            Token(TokenType.IDENT, 'x'),
            Token(TokenType.PLUS, '+'),
            Token(TokenType.IDENT, 'y'),
            Token(TokenType.SEMICOLON, ';'),
            Token(TokenType.RBRACE, '}'),
            Token(TokenType.SEMICOLON, ';'),
        ]

        self.assertEquals(tokens, expected_tokens)

    def test_function_call(self) -> None:
        """Test functino call tokens."""
        
        source: str = 'var resultado = suma(dos, tres);'
        lexer: Lexer = Lexer(source)
        tokens: List[Token] = []

        for i in range(10):
            tokens.append(lexer.next_token())

        expected_tokens: List[Token] = [
            Token(TokenType.LET, 'var'),
            Token(TokenType.IDENT, 'resultado'),
            Token(TokenType.ASSING, '='),
            Token(TokenType.IDENT, 'suma'),
            Token(TokenType.LPAREN, '('),
            Token(TokenType.IDENT, 'dos'),
            Token(TokenType.COMMA, ','),
            Token(TokenType.IDENT, 'tres'),
            Token(TokenType.RPAREN, ')'),
            Token(TokenType.SEMICOLON, ';'),
        ]

        self.assertEquals(tokens, expected_tokens)
