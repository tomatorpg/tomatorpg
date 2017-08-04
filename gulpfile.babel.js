import gulp, { parallel, series } from 'gulp';
import gulpHelp from 'gulp-help-four';
import htmlmin from 'gulp-htmlmin';
import imagemin from 'gulp-imagemin';

gulpHelp(gulp);

gulp.task('build:images', 'build / handle asset images', () =>
  gulp.src('assets/src/images/**/*')
    .pipe(imagemin())
    .pipe(gulp.dest('assets/dist/images')));

gulp.task('build:html', 'minify html template for assets', () =>
  gulp.src('assets/src/html/**/*.html')
    .pipe(htmlmin({ collapseWhitespace: true }))
    .pipe(gulp.dest('assets/dist/html')));

gulp.task('default', false, series('build:html'));

gulp.task('build', false, parallel('build:html'));
