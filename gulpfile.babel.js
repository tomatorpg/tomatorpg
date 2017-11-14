import gulp, { parallel, series } from 'gulp';
import gulpHelp from 'gulp-help-four';
import htmlmin from 'gulp-htmlmin';

gulpHelp(gulp);

gulp.task('build:html', 'minify html template for assets', () =>
  gulp.src('src/html/**/*.html')
    .pipe(htmlmin({ collapseWhitespace: true }))
    .pipe(gulp.dest('assets/dist/html')));

gulp.task('build:images', 'copy images from src to assets/dist', () =>
  gulp.src('src/images/**/*')
    .pipe(gulp.dest('assets/dist/images')));

gulp.task('default', false, series('build:html'));

gulp.task('build', false, parallel('build:html'));
