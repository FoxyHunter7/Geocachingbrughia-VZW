<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('event_languages', function (Blueprint $table) {
            $table->bigIncrements('id');

            $table->foreignId('event_id')->constrained()->onUpdate('cascade')->onDelete('cascade');
            $table->foreignId('lang_code')->constrained(table: 'languages', column: 'code')->onUpdate('cascade')->onDelete('cascade');
            $table->unique(['event_id','lang_code']);

            $table->string('description');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('event_languages');
    }
};
