import {EditorView} from '@codemirror/view';
import {type TagStyle} from '@codemirror/language';
import {type Extension} from '@codemirror/state';

import {
	HighlightStyle,
	syntaxHighlighting,
} from '@codemirror/language';

// Credit: https://github.com/vadimdemedes/thememirror/blob/main/source/create-theme.ts
interface Options {
	/**
	 * Theme variant. Determines which styles CodeMirror will apply by default.
	 */
	variant: Variant;

	/**
	 * Settings to customize the look of the editor, like background, gutter, selection and others.
	 */
	settings: Settings;

	/**
	 * Syntax highlighting styles.
	 */
	styles: TagStyle[];
}

type Variant = 'light' | 'dark';

interface Settings {
	/**
	 * Editor background.
	 */
	background: string;

	/**
	 * Default text color.
	 */
	foreground: string;

	/**
	 * Caret color.
	 */
	caret: string;

	/**
	 * Selection background.
	 */
	selection: string;

	/**
	 * Background of highlighted lines.
	 */
	lineHighlight: string;

	/**
	 * Gutter background.
	 */
	gutterBackground: string;

	/**
	 * Text color inside gutter.
	 */
	gutterForeground: string;
}

const createTheme = ({variant, settings, styles}: Options): Extension => {
	const theme = EditorView.theme(
		{
			 
			'&': {
				backgroundColor: settings.background,
				color: settings.foreground,
				borderRadius: "10px",
				width: "100%",
				height: "83vh",
				border: "1px solid var(--light-light-grey)",
				boxShadow: "0px 0px 35px 5px #0b0b0b",
				fontSize: "1.2em"
			},
			'.cm-content': {
				caretColor: settings.caret,
			},
			'.cm-cursor, .cm-dropCursor': {
				borderLeftColor: settings.caret,
			},
			'&.cm-focused .cm-selectionBackgroundm .cm-selectionBackground, .cm-content ::selection':
				{
					backgroundColor: settings.selection,
				},
			'.cm-activeLine': {
				backgroundColor: settings.lineHighlight,
			},
			'.cm-gutters': {
				backgroundColor: settings.gutterBackground,
				color: settings.gutterForeground,
				borderRadius: "10px 0 0 10px"
			},
			'.cm-activeLineGutter': {
				backgroundColor: settings.lineHighlight,
			},
		},
		{
			dark: variant === 'dark',
		},
	);

	const highlightStyle = HighlightStyle.define(styles);
	const extension = [theme, syntaxHighlighting(highlightStyle)];

	return extension;
};

export default createTheme;
